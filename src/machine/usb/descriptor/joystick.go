package descriptor

var deviceJoystick = [deviceTypeLen]byte{
	deviceTypeLen,
	deviceTypeID,
	0x00, 0x02, // USB version
	0xef,       // device class
	0x02,       // subclass
	0x01,       // protocol
	0x40,       // maxpacketsize
	0x41, 0x23, // vendor id
	0x36, 0x80, // product id
	0x00, 0x01, // device
	0x01, // manufacturer
	0x02, // product
	0x03, // SerialNumber
	0x01, // NumConfigurations
}

var DeviceJoystick = DeviceType{
	data: deviceJoystick[:],
}

var configurationCDCJoystick = [configurationTypeLen]byte{
	configurationTypeLen,
	configurationTypeID,
	0x6b, 0x00, // adjust length as needed
	0x03, // number of interfaces
	0x01, // configuration value
	0x00, // index to string description
	0xa0, // attributes
	0xfa, // maxpower
}

var ConfigurationCDCJoystick = ConfigurationType{
	data: configurationCDCJoystick[:],
}

var interfaceHIDJoystick = [interfaceTypeLen]byte{
	interfaceTypeLen,
	interfaceTypeID,
	0x02, // InterfaceNumber
	0x00, // AlternateSetting
	0x02, // NumEndpoints
	0x03, // InterfaceClass
	0x00, // InterfaceSubClass
	0x00, // InterfaceProtocol
	0x00, // Interface
}

var InterfaceHIDJoystick = InterfaceType{
	data: interfaceHIDJoystick[:],
}

var classHIDJoystick = [classHIDTypeLen]byte{
	classHIDTypeLen,
	classHIDTypeID,
	0x11, // HID version L
	0x01, // HID version H
	0x00, // CountryCode
	0x01, // NumDescriptors
	0x22, // ClassType
	0x00, // ClassLength L
	0x00, // ClassLength H
}

var ClassHIDJoystick = ClassHIDType{
	data: classHIDJoystick[:],
}

var CDCJoystick = Descriptor{
	Device: DeviceJoystick.Bytes(),
	Configuration: appendSlices([][]byte{
		ConfigurationCDCJoystick.Bytes(),
		InterfaceAssociationCDC.Bytes(),
		InterfaceCDCControl.Bytes(),
		ClassSpecificCDCHeader.Bytes(),
		ClassSpecificCDCACM.Bytes(),
		ClassSpecificCDCUnion.Bytes(),
		EndpointEP1IN.Bytes(),
		InterfaceCDCData.Bytes(),
		EndpointEP2OUT.Bytes(),
		EndpointEP3IN.Bytes(),
		InterfaceHIDJoystick.Bytes(),
		ClassHIDJoystick.Bytes(),
		EndpointEP4IN.Bytes(),
		EndpointEP5OUT.Bytes(),
	}),
	HID: map[uint16][]byte{},
}
