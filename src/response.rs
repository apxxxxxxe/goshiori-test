use winapi::{
    shared::{
        minwindef::LPBOOL,
        ntdef::{LPCSTR, NULL},
    },
    um::{stringapiset::WideCharToMultiByte, winnt::LPSTR},
};

use shiorust::message::{parts::HeaderName, parts::*, traits::*, Response};

const CRLF: &str = "\r\n";

pub enum ResponseError {
    DecodeFailed,
}

pub struct PluginResponse {
    pub response: Response,
}

impl std::fmt::Display for PluginResponse {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(
            f,
            "PLUGIN/{} {}{}{}{}",
            self.response.version, self.response.status, CRLF, self.response.headers, CRLF
        )
    }
}

impl PluginResponse {
    pub fn new() -> PluginResponse {
        let mut headers = Headers::new();
        headers.insert(
            HeaderName::Standard(StandardHeaderName::Charset),
            String::from("UTF-8"),
        );

        PluginResponse {
            response: Response {
                version: Version::V20,
                status: Status::OK,
                headers,
            },
        }
    }

    pub fn new_nocontent() -> PluginResponse {
        let mut r = PluginResponse::new();
        r.response.status = Status::NoContent;
        r
    }

    /// 自身をエンコードされた文字バイト列にして返す
    pub fn to_encoded_bytes(&mut self) -> Result<Vec<i8>, ResponseError> {
        let req = self.to_string();

        let mut wide_chars: Vec<u16> = req.encode_utf16().collect();

        const UTF8: u32 = 65001;
        let result = wide_char_to_multi_byte(&mut wide_chars, UTF8)
            .map_err(|_| ResponseError::DecodeFailed)?;

        Ok(result)
    }
}

fn wide_char_to_multi_byte(from: &mut Vec<u16>, codepage: u32) -> Result<Vec<i8>, ()> {
    from.push(0);

    let to_buf_size = unsafe {
        WideCharToMultiByte(
            codepage,
            0,
            from.as_ptr(),
            -1,
            NULL as LPSTR,
            0,
            NULL as LPCSTR,
            NULL as LPBOOL,
        )
    };

    if to_buf_size == 0 {
        return Err(());
    }

    let mut to_buf: Vec<i8> = vec![0; to_buf_size as usize + 1];
    let result = unsafe {
        WideCharToMultiByte(
            codepage,
            0,
            from.as_ptr(),
            -1,
            to_buf.as_mut_ptr(),
            to_buf_size,
            NULL as LPCSTR,
            NULL as LPBOOL,
        )
    };

    if result == 0 {
        Err(())
    } else {
        Ok(to_buf)
    }
}
