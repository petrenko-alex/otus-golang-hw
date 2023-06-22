package hw03frequencyanalysis_test

import (
	"errors"
	hw03frequencyanalysis "github.com/petrenko-alex/otus-golang-hw/hw03_frequency_analysis"
	"testing"

	"github.com/stretchr/testify/require"
)

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

var commonTestCases = []struct {
	desc   string
	input  string
	output []string
}{
	{
		desc:   "empty string",
		input:  "",
		output: []string{},
	},
	{
		desc:   "simple digit string",
		input:  "10 10 10 10 10 10 10 10 10 10 9 9 9 9 9 9 9 9 9 8 8 8 8 8 8 8 8 7 7 7 7 7 7 7 6 6 6 6 6 6 5 5 5 5 5 4 4 4 4 3 3 3 2 2 1",
		output: []string{"10", "9", "8", "7", "6", "5", "4", "3", "2", "1"},
	},
	{
		desc:   "simple digit string, not enough for top 10",
		input:  "3 3 3 2 2 1",
		output: []string{"3", "2", "1"},
	},
	{
		desc:   "simple word string",
		input:  "cat dog bird dog cat cat cat",
		output: []string{"cat", "dog", "bird"},
	},
	{
		desc:   "mixed digit word string",
		input:  "3 3 cat 3 cat 2 dog 2 cat",
		output: []string{"3", "cat", "2", "dog"},
	},
	{
		desc:   "one word string",
		input:  "cat",
		output: []string{"cat"},
	},
	{
		desc:   "one word repeated",
		input:  "cat cat cat cat",
		output: []string{"cat"},
	},
	{
		desc:   "equal frequency, sorting",
		input:  "man cat bird dog",
		output: []string{"bird", "cat", "dog", "man"},
	},
	{
		desc:   "cyrillic string",
		input:  "кошка собака птица кошка кошка собака",
		output: []string{"кошка", "собака", "птица"},
	},
	{
		desc:   "word form",
		input:  "кошка кошкой кошка кошке кошкой кошка",
		output: []string{"кошка", "кошкой", "кошке"},
	},
	{
		desc:   "more than top 10",
		input:  "a a b b c c d d e e f f g g h h i i j j k k l l m m n o p",
		output: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
	},
	{
		desc:   "symbols and emojis",
		input:  "😀 a 😀 🤣 a a",
		output: []string{"a", "😀", "🤣"},
	},
	{
		desc:   "different whitespaces",
		input:  "a  a		b\tc",
		output: []string{"a", "b", "c"},
	},
}

func TestTop10PositivePunctuation(t *testing.T) {
	testCases := []struct {
		desc   string
		input  string
		output []string
	}{
		{
			desc:   "capital letters",
			input:  "Cat dog Cat cat cat Cat",
			output: []string{"cat", "dog"},
		},
		{
			desc:   "punctuation, commas",
			input:  "cat and dog, one dog, two cats and one man",
			output: []string{"and", "dog", "one", "cat", "cats", "man", "two"},
		},
		{
			desc:   "punctuation, dash",
			input:  "cat - dog cat. man -",
			output: []string{"cat", "dog", "man"},
		},
		{
			desc:  "complex text",
			input: text,
			output: []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			},
		},
	}
	testCases = append(commonTestCases, testCases...)
	top10 := hw03frequencyanalysis.GeneralTextWordFrequency{
		TextValidator:    hw03frequencyanalysis.Utf8Validator{},
		FrequencyCounter: hw03frequencyanalysis.PunctuationFrequencyCounter{},
		FrequencySorter:  hw03frequencyanalysis.DescendingFrequencySorter{},
		FrequencyLimiter: hw03frequencyanalysis.SimpleFrequencyLimiter{Limit: 10},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.desc, func(t *testing.T) {
			res, _ := top10.Top(testCase.input)

			require.Equal(t, testCase.output, res)
		})
	}
}

func TestTop10PositiveNonPunctuation(t *testing.T) {
	testCases := []struct {
		desc   string
		input  string
		output []string
	}{
		{
			desc:   "capital letters",
			input:  "Cat dog Cat cat cat Cat",
			output: []string{"Cat", "cat", "dog"},
		},
		{
			desc:   "punctuation, commas",
			input:  "cat and dog, one dog,two cats and one man",
			output: []string{"and", "one", "cat", "cats", "dog,", "dog,two", "man"},
		},
		{
			desc:   "punctuation, dash",
			input:  "cat - dog cat. man -",
			output: []string{"-", "cat", "cat.", "dog", "man"},
		},
		{
			desc:  "complex text",
			input: text,
			output: []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			},
		},
	}
	testCases = append(commonTestCases, testCases...)
	top10 := hw03frequencyanalysis.GeneralTextWordFrequency{
		TextValidator:    hw03frequencyanalysis.Utf8Validator{},
		FrequencyCounter: hw03frequencyanalysis.NonPunctuationFrequencyCounter{},
		FrequencySorter:  hw03frequencyanalysis.DescendingFrequencySorter{},
		FrequencyLimiter: hw03frequencyanalysis.SimpleFrequencyLimiter{Limit: 10},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.desc, func(t *testing.T) {
			res, _ := top10.Top(testCase.input)

			require.Equal(t, testCase.output, res)
		})
	}
}

func TestTop10Errors(t *testing.T) {
	testCases := []struct {
		desc          string
		input         string
		executor      hw03frequencyanalysis.TextWordFrequency
		expectedError error
	}{
		{
			desc:  "PunctuationFrequencyCounter: Invalid UTF-8",
			input: "\xe0 \xe1 \xe2 \xe3 \xe9",
			executor: hw03frequencyanalysis.GeneralTextWordFrequency{
				TextValidator:    hw03frequencyanalysis.Utf8Validator{},
				FrequencyCounter: hw03frequencyanalysis.PunctuationFrequencyCounter{},
				FrequencySorter:  hw03frequencyanalysis.DescendingFrequencySorter{},
				FrequencyLimiter: hw03frequencyanalysis.SimpleFrequencyLimiter{Limit: 10},
			},
			expectedError: hw03frequencyanalysis.InvalidUtf8TextError,
		},
		{
			desc:  "NonPunctuationFrequencyCounter: Invalid UTF-8",
			input: "\xe0 \xe1 \xe2 \xe3 \xe9",
			executor: hw03frequencyanalysis.GeneralTextWordFrequency{
				TextValidator:    hw03frequencyanalysis.Utf8Validator{},
				FrequencyCounter: hw03frequencyanalysis.NonPunctuationFrequencyCounter{},
				FrequencySorter:  hw03frequencyanalysis.DescendingFrequencySorter{},
				FrequencyLimiter: hw03frequencyanalysis.SimpleFrequencyLimiter{Limit: 10},
			},
			expectedError: hw03frequencyanalysis.InvalidUtf8TextError,
		},
	}

	for i, _ := range testCases {
		testCase := testCases[i]
		t.Run(testCase.desc, func(t *testing.T) {
			_, err := testCase.executor.Top(testCase.input)

			require.Truef(t, errors.Is(err, testCase.expectedError), "actual error %q", err)
		})
	}
}
