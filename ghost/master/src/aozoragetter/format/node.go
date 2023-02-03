package format

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type TargetNode struct {
	Element string
	Attr    *html.Attribute
}

func findMainText(node *html.Node) *html.Node {
	var results []*html.Node

	targetTitle := TargetNode{
		Element: "h1",
		Attr: &html.Attribute{
			Key: "class",
			Val: "title",
		},
	}
	results = FindNode(node, targetTitle, []*html.Node{})
	if len(results) != 1 {
		return nil
	}
	title := results[0]
	title.Parent.RemoveChild(title)

	targetAuthor := TargetNode{
		Element: "h2",
		Attr: &html.Attribute{
			Key: "class",
			Val: "author",
		},
	}
	results = FindNode(node, targetAuthor, []*html.Node{})
	if len(results) != 1 {
		return nil
	}
	author := results[0]
	author.Parent.RemoveChild(author)

	targetMain := TargetNode{
		Element: "div",
		Attr: &html.Attribute{
			Key: "class",
			Val: "main_text",
		},
	}
	results = FindNode(node, targetMain, []*html.Node{})
	if len(results) != 1 {
		return nil
	}
	mainText := results[0]

	br := make([]*html.Node, 3)
	for i := range br {
		br[i] = &html.Node{
			Type:      html.ElementNode,
			Data:      "br",
			DataAtom:  atom.Br,
			Attr:      []html.Attribute{},
			Namespace: "",
		}

	}

	mainText.InsertBefore(br[2], mainText.FirstChild)
	mainText.InsertBefore(br[1], mainText.FirstChild)
	mainText.InsertBefore(author, mainText.FirstChild)
	mainText.InsertBefore(br[0], mainText.FirstChild)
	mainText.InsertBefore(title, mainText.FirstChild)

	return mainText
}

func FindNode(node *html.Node, target TargetNode, result []*html.Node) []*html.Node {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == target.Element {
			if c.Attr != nil {
				for _, a := range c.Attr {
					if a.Key == target.Attr.Key && a.Val == target.Attr.Val {
						result = append(result, c)
					}
				}
			} else {
				result = append(result, c)
			}
		}
		result = append(result, FindNode(c, target, []*html.Node{})...)
	}
	return result
}

func extractClass(node *html.Node) string {
	return extractAttr(node, "class")
}

func extractAttr(node *html.Node, attr string) string {
	for _, a := range node.Attr {
		if a.Key == attr {
			return a.Val
		}
	}
	return ""
}
