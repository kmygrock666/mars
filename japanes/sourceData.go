package main

import (
	"math/rand"
	"time"
)

//平假名
var flagAccent = []string{
	"あ", "い", "う", "え", "お", "ん",
	"か", "き", "く", "け", "こ", "が", "ぎ", "ぐ", "げ", "ご",
	"さ", "し", "す", "せ", "そ", "ざ", "じ", "ず", "ぜ", "ぞ",
	"た", "ち", "つ", "て", "と", "だ", "ぢ", "づ", "で", "ど",
	"な", "に", "ぬ", "ね", "の",
	"は", "ち", "つ", "て", "と", "ば", "び", "ぶ", "べ", "ぼ", "ぱ", "ぴ", "ぷ", "ぺ", "ぽ",
	"ま", "ち", "つ", "て", "と",
	"ら", "ち", "つ", "て", "と",
	"や", "ゆ", "よ", "わ", "を",
}

//平假名(字源)
var flagWord = []string{
	"安", "以", "宇", "衣", "於", "无",
	"加", "幾", "久", "計", "己", "加", "幾", "久", "計", "己",
	"左", "之", "寸", "世", "曾", "左", "之", "寸", "世", "曾",
	"太", "知", "川", "天", "止", "太", "知", "川", "天", "止",
	"奈", "仁", "奴", "祢", "乃",
	"波", "比", "不", "部", "保", "波", "比", "不", "部", "保", "波", "比", "不", "部", "保",
	"末", "美", "武", "女", "毛",
	"良", "利", "留", "礼", "呂",
	"也", "由", "与", "和", "无",
}

//平假名(羅馬拼音)
var flagPinyin = []string{
	"a", "i", "u", "e", "o", "n",
	"ka", "ki", "ku", "ke", "ko", "ga", "gi", "gu", "ge", "go",
	"sa", "si/shi", "su", "se", "so", "za", "zi", "zu", "ze", "zo",
	"ta", "chi", "tsu", "te", "to", "da", "zi", "zu", "de", "do",
	"na", "ni", "nu", "ne", "no",
	"ha", "hi", "hu/fu", "he", "ho", "ba", "bi", "bu", "be", "bo", "pa", "pi", "pu", "pe", "po",
	"ma", "mi", "mu", "me", "mo",
	"ra", "ri", "ru", "re", "ro",
	"ya", "yu", "yo", "wa", "wo",
}

//片假名
var sliceAccent = []string{
	"ア", "イ", "ウ", "エ", "オ", "ン",
	"カ", "キ", "ク", "ケ", "コ", "ガ", "ギ", "グ", "ゲ", "ゴ",
	"サ", "シ", "ス", "セ", "ソ", "ザ", "ジ", "ズ", "ゼ", "ゾ",
	"タ", "チ", "ツ", "テ", "ト", "ダ", "ヂ", "ヅ", "デ", "ド",
	"ナ", "ニ", "ヌ", "ネ", "ノ",
	"ハ", "ヒ", "フ", "ヘ", "ホ", "バ", "ビ", "ブ", "ベ", "ボ", "パ", "ピ", "プ", "ペ", "ポ",
	"マ", "ミ", "ム", "メ", "モ",
	"ラ", "リ", "ル", "レ", "ロ",
	"ヤ", "ユ", "ヨ", "ワ", "ヲ",
}

//片假名(字源)
var sliceWord = []string{
	"阿", "伊", "宇", "江", "於", "尔",
	"加", "幾", "久", "介", "己", "加", "幾", "久", "介", "己",
	"散", "之", "須", "世", "曽", "散", "之", "須", "世", "曽",
	"多", "千", "川", "天", "止", "多", "千", "川", "天", "止",
	"奈", "仁", "奴", "祢", "乃",
	"八", "比", "不", "部", "保", "八", "比", "不", "部", "保", "八", "比", "不", "部", "保",
	"末", "三", "牟", "女", "毛",
	"良", "利", "留", "礼", "呂",
	"也", "由", "與", "和", "无",
}

//片假名(羅馬拼音)
var slicePinyin = []string{
	"a", "i", "u", "e", "o", "n",
	"ka", "ki", "ku", "ke", "ko", "ga", "gi", "gu", "ge", "go",
	"sa", "si/shi", "su", "se", "so", "za", "zi", "zu", "ze", "zo",
	"ta", "chi", "tsu", "te", "to", "da", "zi", "zu", "de", "do",
	"na", "ni", "nu", "ne", "no",
	"ha", "hi", "hu/fu", "he", "ho", "ba", "bi", "bu", "be", "bo", "pa", "pi", "pu", "pe", "po",
	"ma", "mi", "mu", "me", "mo",
	"ra", "ri", "ru", "re", "ro",
	"ya", "yu", "yo", "wa", "wo",
}

var tmpIndex []int

//JapanAccent JapanAccent
type JapanAccent struct {
	Foreign string
	Native  string
	Pinyin  string
}

func newRoundshuffle() {
	tmpIndex = nil
	var count = len(flagAccent) + len(sliceAccent)
	for i := 0; i < count; i++ {
		tmpIndex = append(tmpIndex, i)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(tmpIndex), func(i, j int) { tmpIndex[i], tmpIndex[j] = tmpIndex[j], tmpIndex[i] })
}

func getFlagAccentWord(fIndex int) *JapanAccent {
	return &JapanAccent{
		Foreign: sliceAccent[fIndex],
		Native:  sliceWord[fIndex],
		Pinyin:  slicePinyin[fIndex],
	}
}

func getSliceAccentWord(sIndex int) *JapanAccent {
	return &JapanAccent{
		Foreign: flagAccent[sIndex],
		Native:  flagWord[sIndex],
		Pinyin:  flagPinyin[sIndex],
	}
}

func getAccentWord() *JapanAccent {
	i := tmpIndex[0]
	var ret *JapanAccent
	if i > 70 {
		i = i - 71
		ret = getSliceAccentWord(i)
	} else {
		ret = getFlagAccentWord(i)
	}

	return ret
}

func (ja *JapanAccent) getCurrentAccent(wordIdx int64) {

}

//CheckPinyin CheckPinyin
func (ja *JapanAccent) CheckPinyin(pinyin string) bool {
	ret := false
	_ = getAccentWord()
	if ja.Pinyin == pinyin {
		ret = true
		tmpIndex = tmpIndex[1:]
	}

	return ret
}
