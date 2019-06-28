package xml_lib

import "encoding/xml"

func XmlToObject(strXML string, v interface{}) (err error) {
	err = xml.Unmarshal([]byte(strXML), &v)
	return
}

func ObjectToXml(v interface{}) (strXml string, err error) {
	var b []byte
	b, err = xml.Marshal(v)
	if err == nil {
		strXml = string(b)
	}
	return
}

func ObjectToObject(desc, source interface{}) (err error) {
	var strXml string
	strXml, err = ObjectToXml(source)
	if err == nil {
		err = XmlToObject(strXml, &desc)
	}
	return
}
