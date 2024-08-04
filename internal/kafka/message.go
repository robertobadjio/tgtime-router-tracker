package kafka

type InOfficeMessage struct {
	MacAddress string
}

const InOfficeTopic = "in-office"
const partition = 0
