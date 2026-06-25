package modelschema

import "encoding/xml"

// mo phong cau truc the <get-config> của NETCONF
type GetConfigRequest struct {
	XMLName xml.Name `xml:"urn.ietf:params:xml:ns:netconf:base:1.0 get-config"`
	Source  Source   `xml:"source"`
	Filter  Filter   `xml:"filter"`
}

type Source struct {
	Running struct{} `xml:"running"`
}

type Filter struct {
	Type  string `xml:"type,attr"` //attr type="subtree"
	Inner string `xml:",innerxml"` //Inject RawXML(dataPayload) without encoding
}

// RPCReply struct contains all content from tag <rpc-reply>
type RPCReply struct {
	InnerXML string `xml:",innerxml"`
}

// //Marshal method interface for netconf.Exec
// type RawXML string
// func (r RawXML) MarshalMethod() string {
// 	return string(r)
// }
