mod client;
mod events;
mod request;
mod response;

use client::RpcClient;
use request::PluginRequest;

use std::fs::File;
use std::path::Path;
use std::slice;

use once_cell::sync::Lazy;

use winapi::ctypes::c_long;
use winapi::shared::minwindef::{BOOL, DWORD, HGLOBAL, HINSTANCE, LPVOID, MAX_PATH, TRUE};
use winapi::um::libloaderapi::GetModuleFileNameW;
use winapi::um::winbase::{GlobalAlloc, GlobalFree, GMEM_FIXED};
use winapi::um::winnt::{
    DLL_PROCESS_ATTACH, DLL_PROCESS_DETACH, DLL_THREAD_ATTACH, DLL_THREAD_DETACH,
};

use shiorust::message::Parser;

#[macro_use]
extern crate log;
extern crate simplelog;

use simplelog::*;

static mut DLL_PATH: String = String::new();
static mut RPCCLIENT: Lazy<RpcClient> = Lazy::new(|| RpcClient::new());

#[no_mangle]
pub extern "system" fn DllMain(
    h_module: HINSTANCE,
    ul_reason_for_call: DWORD,
    _l_reserved: LPVOID,
) -> BOOL {
    match ul_reason_for_call {
        DLL_PROCESS_ATTACH => {
            debug!("DLL_PROCESS_ATTACH");
            register_dll_path(h_module);
            let path;
            unsafe {
                path = Path::new(&DLL_PATH.clone())
                    .parent()
                    .unwrap()
                    .join("ukaing.log");
            };
            WriteLogger::init(
                LevelFilter::Debug,
                Config::default(),
                File::create(path).unwrap(),
            )
            .unwrap();
        }
        DLL_PROCESS_DETACH => {
            debug!("DLL_PROCESS_DETACH");
        }
        DLL_THREAD_ATTACH => {}
        DLL_THREAD_DETACH => {
            debug!("DLL_THREAD_DETACH");
        }
        _ => {}
    }
    return TRUE;
}

fn register_dll_path(h_module: HINSTANCE) {
    let mut buf: [u16; MAX_PATH + 1] = [0; MAX_PATH + 1];
    unsafe {
        GetModuleFileNameW(h_module, buf.as_mut_ptr(), MAX_PATH as u32);
    }

    let p = buf.partition_point(|v| *v != 0);

    unsafe {
        DLL_PATH = String::from_utf16_lossy(&buf[..p]);
    }
}

#[no_mangle]
pub extern "cdecl" fn load(h: HGLOBAL, _len: c_long) -> BOOL {
    unsafe { GlobalFree(h) };

    debug!("load");
    unsafe { RPCCLIENT.start() };

    return TRUE;
}

#[no_mangle]
pub extern "cdecl" fn unload() -> BOOL {
    debug!("unload");
    unsafe { RPCCLIENT.close() };
    return TRUE;
}

#[no_mangle]
pub extern "cdecl" fn request(h: HGLOBAL, len: *mut c_long) -> HGLOBAL {
    // リクエストの取得
    let v = unsafe { hglobal_to_vec_u8(h, *len) };
    unsafe { GlobalFree(h) };

    let s = String::from_utf8_lossy(&v).to_string();

    let pr = PluginRequest::parse(&s).unwrap();
    let r = pr.request;

    let mut response;
    response = events::handle_request(&r);

    let response_bytes = response.to_encoded_bytes().unwrap_or(Vec::new());

    let h = slice_i8_to_hglobal(len, &response_bytes);

    h
}

fn slice_i8_to_hglobal(h_len: *mut c_long, data: &[i8]) -> HGLOBAL {
    let data_len = data.len();

    let h = unsafe { GlobalAlloc(GMEM_FIXED, data_len) };

    unsafe { *h_len = data_len as c_long };

    let h_slice = unsafe { slice::from_raw_parts_mut(h as *mut i8, data_len) };

    for (index, value) in data.iter().enumerate() {
        h_slice[index] = *value;
    }

    return h;
}

fn hglobal_to_vec_u8(h: HGLOBAL, len: c_long) -> Vec<u8> {
    let mut s = vec![0; len as usize + 1];

    let slice = unsafe { slice::from_raw_parts(h as *const u8, len as usize) };

    for (index, value) in slice.iter().enumerate() {
        s[index] = *value;
    }
    s[len as usize] = b'\0';

    return s;
}

#[cfg(test)]
mod test {
    const TOKEN: &str = "1033946714102562826";
    use discord_rich_presence::{
        activity::{Activity, Button, Timestamps},
        DiscordIpc, DiscordIpcClient,
    };
    use std::time::Duration;

    #[test]
    fn test_client() {
        let mut client = DiscordIpcClient::new(TOKEN).unwrap();
        client.connect().unwrap();
        let activity = Activity::new()
            .state("t ")
            .timestamps(Timestamps::new().start(chrono::Local::now().timestamp()))
            .buttons(vec![Button::new(
                "配布元 / craftmanurl",
                "https://anarchytansu.pv.land.to/",
            )]);
        client.set_activity(activity).unwrap();
        std::thread::sleep(Duration::from_secs(10));
        client.close().unwrap();
    }
}
