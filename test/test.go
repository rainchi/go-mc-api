package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	nbt2Json "github.com/Lirsty/nbt2json"
	nbtTool "github.com/rain931215/go-mc-api/nbt"
	t "github.com/rain931215/go-mc-api/plugin/autobuilder"
)

func main() {
	//readGzip("TreeFarm_by_Ian0822.litematic")
	writeLitematic()
}

func writeLitematic() {
	l := t.NewLitematic()
	l.Metadata.Auther = "Lirsty"
	l.Metadata.Description = "none"
	l.Metadata.EnclosingSize.X = 1
	l.Metadata.EnclosingSize.Y = 2
	l.Metadata.EnclosingSize.Z = 3
	l.Metadata.Name = "test"
	l.Metadata.RegionCount = 2
	l.Metadata.TimeCreated = 123456789
	l.Metadata.TimeModified = 12345678910
	l.Metadata.TotalBlocks = 101
	l.Metadata.TotalVolume = 1011

	l.Regions.Regions = make(map[string]t.Region)
	r := new(t.Region)
	r.Position.X = 1
	r.Position.Y = 12
	r.Position.Z = 123
	r.Size.X = 321
	r.Size.Y = 32
	r.Size.Z = 3
	air := new(t.Blocktype)
	air.Properties = make(map[string]string)
	air.Name = "air"
	air.Properties["AirTest"] = "Hello"
	air.Properties["AirTest2"] = "Hello"
	dirt := new(t.Blocktype)
	dirt.Properties = make(map[string]string)
	dirt.Name = "dirt"
	dirt.Properties["DirtTest"] = "Hello"
	r.BlockStatePalette.Blocks = []t.Blocktype{*air, *dirt}
	r.BlockStates = []int64{123, 123456, 12345678}
	l.Regions.Regions["region01"] = *r
	l.Regions.Regions["region02"] = *r
	l.WriteFile("test.nbt")
}

func bigTest() {
	nbt2Json.UseJavaEncoding()
	//nbt2Json.UseLongAsString()
	nbt := nbtTool.NewNbt()
	Level := nbtTool.NewCompoundTag("Level")
	nbt.AddCompoundTag(Level)

	nestedCompoundTest := nbtTool.NewCompoundTag("nestedCompoundTest")
	egg := nbtTool.NewCompoundTag("egg")
	egg.AddNewValue("name", "Eggbert")
	egg.AddNewValue("value", float32(0.5))
	ham := nbtTool.NewCompoundTag("ham")
	ham.AddNewValue("name", "hampus")
	ham.AddNewValue("value", float32(0.75))
	nestedCompoundTest.AddCompoundTag(egg)
	nestedCompoundTest.AddCompoundTag(ham)
	Level.AddCompoundTag(nestedCompoundTest)

	Level.AddNewValue("intTest", int32(2147483647))
	Level.AddNewValue("byteTest", byte(127))
	Level.AddNewValue("stringTest", "HELLO WORLD THIS IS A TEST STRING")

	listTestLong := nbtTool.NewListTag("listTestLong", nbtTool.TagLong)
	listTestLong.AddNewValue(int64(11))
	listTestLong.AddNewValue(int64(12))
	listTestLong.AddNewValue(int64(13))
	listTestLong.AddNewValue(int64(14))
	listTestLong.AddNewValue(int64(15))
	Level.AddListTag(listTestLong)

	Level.AddNewValue("doubleTest", float64(0.49312871321823148))
	Level.AddNewValue("floatTest", float64(0.49823147058486938))
	Level.AddNewValue("longTest", int64(9223372036854775807))

	listTestCompound := nbtTool.NewListTag("listTestCompound", nbtTool.TagCompound)
	none0 := nbtTool.NewCompoundTag("none0")
	none0.AddNewValue("created-on", int64(1264099775885))
	none0.AddNewValue("name", "Compound tag #0")
	none1 := nbtTool.NewCompoundTag("none1")
	none1.AddNewValue("created-on", int64(1264099775885))
	none1.AddNewValue("name", "Compound tag #1")
	listTestCompound.AddCompoundTag(none0)
	listTestCompound.AddCompoundTag(none1)
	Level.AddListTag(listTestCompound)

	bytes := []byte{}
	for n := 0; n <= 1000; n++ {
		bytes = append(bytes, byte((n*n*255+n*7)%100))
	}
	Level.AddNewValue(`byteArrayTest (the first 1000 values of (n*n*255+n*7)%100, starting with n=0 (0, 62, 34, 16, 8, ...))`, bytes)
	Level.AddNewValue("shortTest", int16(32767))
	int64s := []int64{9223372036854775807, 15151515}
	Level.AddNewValue("LongArrayTest", int64s)

	data, err := nbt.ToJson()
	checkerr(err)
	fmt.Println(string(data))

	out, err := nbt2Json.Json2Nbt(data)
	checkerr(err)
	fmt.Println(out)
	err = ioutil.WriteFile("test.nbt", out, 0622)
	checkerr(err)

	read, err := ioutil.ReadFile("test.nbt")
	checkerr(err)

	data, err = nbt2Json.Nbt2Json(read, "bigtest")
	checkerr(err)
	fmt.Println(string(data))
}
func readGzip(file string) {
	nbt2Json.UseJavaEncoding()
	read, err := ioutil.ReadFile(file)
	checkerr(err)
	rdata := bytes.NewReader(read)
	r, _ := gzip.NewReader(rdata)
	s, _ := ioutil.ReadAll(r)
	data, err := nbt2Json.Nbt2Json(s, "test")
	checkerr(err)
	fmt.Println(string(data))
}

func checkerr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
