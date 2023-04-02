package descriptor

import (
	"encoding/binary"
)

var configurationCDCHID = [configurationTypeLen]byte{
	configurationTypeLen,
	configurationTypeID,
	0x64, 0x00, // adjust length as needed
	0x03, // number of interfaces
	0x01, // configuration value
	0x00, // index to string description
	0xa0, // attributes
	0x32, // maxpower
}

var ConfigurationCDCHID = ConfigurationType{
	data: configurationCDCHID[:],
}

var interfaceHID = [interfaceTypeLen]byte{
	interfaceTypeLen,
	interfaceTypeID,
	0x02, // InterfaceNumber
	0x00, // AlternateSetting
	0x01, // NumEndpoints
	0x03, // InterfaceClass
	0x00, // InterfaceSubClass
	0x00, // InterfaceProtocol
	0x00, // Interface
}

var InterfaceHID = InterfaceType{
	data: interfaceHID[:],
}

const (
	classHIDTypeLen = 9
	classHIDTypeID  = 0x21
)

type ClassHIDType struct {
	data []byte
}

func (d ClassHIDType) Bytes() []byte {
	return d.data[:]
}

func (d ClassHIDType) Length(v uint8) {
	d.data[0] = byte(v)
}

func (d ClassHIDType) Type(v uint8) {
	d.data[1] = byte(v)
}

func (d ClassHIDType) HID(v uint16) {
	binary.LittleEndian.PutUint16(d.data[2:4], v)
}

func (d ClassHIDType) CountryCode(v uint8) {
	d.data[4] = byte(v)
}

func (d ClassHIDType) NumDescriptors(v uint8) {
	d.data[5] = byte(v)
}

func (d ClassHIDType) ClassType(v uint8) {
	d.data[6] = byte(v)
}

func (d ClassHIDType) ClassLength(v uint16) {
	binary.LittleEndian.PutUint16(d.data[7:9], v)
}

var classHID = [classHIDTypeLen]byte{
	classHIDTypeLen,
	classHIDTypeID,
	0x01, // HID version L
	0x01, // HID version H
	0x00, // CountryCode
	0x01, // NumDescriptors
	0x22, // ClassType
	0x7E, // ClassLength L
	0x00, // ClassLength H
}

var ClassHID = ClassHIDType{
	data: classHID[:],
}

var CDCHID = Descriptor{
	Device: DeviceCDC.Bytes(),
	Configuration: appendSlices([][]byte{
		ConfigurationCDCHID.Bytes(),
		InterfaceAssociationCDC.Bytes(),
		InterfaceCDCControl.Bytes(),
		ClassSpecificCDCHeader.Bytes(),
		ClassSpecificCDCACM.Bytes(),
		ClassSpecificCDCUnion.Bytes(),
		ClassSpecificCDCCallManagement.Bytes(),
		EndpointEP1IN.Bytes(),
		InterfaceCDCData.Bytes(),
		EndpointEP2OUT.Bytes(),
		EndpointEP3IN.Bytes(),
		InterfaceHID.Bytes(),
		ClassHID.Bytes(),
		EndpointEP4IN.Bytes(),
	}),
	HID: map[uint16][]byte{
		2: []byte{
			// keyboard and mouse
			0x05, 0x01, 0x09, 0x06, 0xa1, 0x01, 0x85, 0x02, 0x05, 0x07, 0x19, 0xe0, 0x29, 0xe7, 0x15, 0x00,
			0x25, 0x01, 0x75, 0x01, 0x95, 0x08, 0x81, 0x02, 0x95, 0x01, 0x75, 0x08, 0x81, 0x03, 0x95, 0x06,
			0x75, 0x08, 0x15, 0x00, 0x25, 0xFF, 0x05, 0x07, 0x19, 0x00, 0x29, 0xFF, 0x81, 0x00, 0xc0, 0x05,
			0x01, 0x09, 0x02, 0xa1, 0x01, 0x09, 0x01, 0xa1, 0x00, 0x85, 0x01, 0x05, 0x09, 0x19, 0x01, 0x29,
			0x03, 0x15, 0x00, 0x25, 0x01, 0x95, 0x03, 0x75, 0x01, 0x81, 0x02, 0x95, 0x01, 0x75, 0x05, 0x81,
			0x03, 0x05, 0x01, 0x09, 0x30, 0x09, 0x31, 0x09, 0x38, 0x15, 0x81, 0x25, 0x7f, 0x75, 0x08, 0x95,
			0x03, 0x81, 0x06, 0xc0, 0xc0,

			0x05, 0x0C, //       Usage Page (Consumer)
			0x09, 0x01, //       Usage (Consumer Control)
			0xA1, 0x01, //       Collection (Application)
			0x85, 0x03, //         Report ID (3)
			0x15, 0x00, //         Logical Minimum (0)
			0x26, 0xFF, 0x1F, //   Logical Maximum (8191)
			0x19, 0x00, //         Usage Minimum (Unassigned)
			0x2A, 0xFF, 0x1F, //   Usage Maximum (0x1FFF)
			0x75, 0x10, //         Report Size (16)
			0x95, 0x01, //         Report Count (1)
			0x81, 0x00, //         Input (Data,Array,Abs,No Wrap,Linear,Preferred State,No Null Position)
			0xC0, //             End       Collection
		},
	},
}
