package platform

import "testing"

func TestLengthEncodeDecode(t *testing.T) {
	for number := 0; number < 10000; number++ {
		length_array := getDataLengthInBytes(int32(number))
		decoded_length := getLengthFromBytes(length_array)
		if number != decoded_length {
			t.Errorf("Expected value: %d, Got value %d", number, decoded_length)
		}
	}
}
