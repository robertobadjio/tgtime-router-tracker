package kafka

type InOfficeMessage struct {
	MacAddress string
}

const inOfficeTopic = "in-office"
const partition = 0
