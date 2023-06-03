package bencode_test

import (
	"bytes"
	"fmt"
	"math"
)

var (
	extremelyLongString = "uhuauugdhzazmrolvdicxxuwvurgpzeuhapuoijmszsfjqdtphmbsnalffyfjiwqzzgohyarvdvqxjcggmwphfozqnkdrwrlytvcoyvjgftubrkjvytwtpscoinckjxgrojcznytayroynzezbdpalovboqhxpwhzdyaethclymjqnggapmsceegihlnwwggmrdkzkipmmwnszulpxfuohfqwrpglnpifevwjkxqkvvypfghnpoejwioltpwuvuemdcifbeerzfvjoryiimktkclvmuczikvufgkpczotqyezadluceulevcwfnwvyejnqoglzsatemozuajmdwgukmztahpzrhcfuhhwxkdbjpnogcmtnbwevozdtuchelwzudctzpvmwmextgcrdgucqaabgwapxvvcsbjebujogevwevzzmjgovqumwgexzbuskwqcjgxtusdcrclqzndieqohfeozmxbtaqrnovslakilhbpfycuhccifzxplwlatxxnklweidpcxcxlotxrzjnmnfcpmlzwahipqtokgwfqnsfwrrxdipyoygkhazlacpzeitjnnyiorrvghqmeebxfymgljynxiusvfyfrlqmqopkqkigblihufnnsuhflbtmotkkyebdvcybemxjuksnvdjqvefhgoiinyyfaxfqfyjuuvgdjqqxdeoregckydlxhagnjlxljvzagjkakdybkwdkucqmjcmlpzexcbocegrhjtgbyawgrzkljjvrqvkpflqsmocszccdnsxzqhaslavcjczosthibrnyogdxjjgyojqsautnqhrnsyuchjribqcecstjnubakzdtebduogpnbzqfismlvfrouivdzrxaucpeoocpesbzpucbtsvdfkpfitsnkxskztyoswydpkvobpjqpthpcwhuumymyknsdjycnzplksiheothiwakemcgyykrqhulfybxhhrqpnlhzvweopxuyttrdnzopaqktumewdkxhgzwkhlcafcladmhlrgaztgprkonsryirrawtpawpegztllpftyrseojsrxxrycsydvackuhlkawzzcqrikhcpwghiljgnduycrdmnqhdxgmmdnpnfhgzpebkkfhxghsfacqabwadtlbhejjrpohgxebwzdqlbelcqhotfbfevscpeauatzcictcwkavbbdvjniflfpnmnattefdgxbynirvfdptagdbvwutzuwrcqtjzoqqvpvlautvpukalgfwqsohtgesgfxlrkevrbzlmdrotrnwwnmkkefmjfjlsqmeookbipokxxfxlgdnsojrhanvniidveyntatncdrtqetjpivsfprmywzilmjreuiawotmzomccdydwmujmmzdmvaokjyltqmrpqoshqwwlmefwyujmnxiceexhwisjnxynrjdrsuoxsvhfypzuksllgpjeisnucvsimfvjlgedxcsiwxvvflblecskdbcsfissxxjedwqxehiepfbevcbfecnrenevfhenasvlbndblpwcxiqhanntdpvpkualnrrxszwwqgmpnrjxbyeevqdeowmxhewovcpbcgvadormsvmozphuzmtkfyvmeqqorrazbtyfzohegvlggacyzdzvgnotfbmjwsdltliaimjjcstmvmsubjcsitazzlovkkzbwsafbssiutlskmkcuceikjtjcdijsiodbhsaqvufhibhcvifbusfjuygtgvnxlsmuypqyzdrgpjlbyzbljjozeyhvwbifxsxmmdaywoqygamwydwuevyaixnznkbvoqydduhfkyjzsqekgcsrhgghctwsavcfrpbjnusodxjpndwwjdktzpureiomrtskzekrgabldgwmatvjwcgemefqvzvalzeyqyluvyqpiqjfudibpiafaiuafyezamuzyvsgmqlmeaublbxavuoejpmtsidrykizljpnccycwuflycxkukhruvbedhlywcvqpienowaggxpknqkzmurzbgyduyjbgojolaittdurksyjffvqkocahvtivanbclahcoyjnqcfmkzfughvbzgsighooecdrndjivdcbfvdgfyynpreqvnqrfvfayvoxijfjkocfvmfsqxpawfmhkdthqhvudzxwpdrasmobxbhvraleceiufqnulwrgldmcdowpatieyczexnykqhstnsmkzvlrsdtzieklpmdwvnumctsrgdegbilqihhxslowkxhvfqgrfuovrldlskavotpksqwxbexlrvfqyrtgcenmoxoegwwfoivpotjidrtjnddbfkjtnxjwumxodlbnddhpubsqbibxkmskyjmvbckoccsordfutdyflzejbyrisovejjtyxmqspkwvupgjjhuivkpbxkbkurawflvbyxuwjpiaysqhureokmjfsxondnkfpgohgobcgsyhdmuozvjubjqfgzwyorghftzlcjawqzkpvaysbiiyikfhuxkzlymyhumcqbvhgmrfqcgyhwrnwkitjgmoccuizkqzftwzvymvaguhivwfppyobkrdevizhfpfwlkibxzjvleuijgtakovepqrjpnhkjxlailtxblqeowfxdztwezqalodhonyuotovwdsaqmpyvswdeajckhwzqbggscphylkcqlatauyhllrojgpskqmqrqmtirnbdnmhvlxnlmequjrbrgbocytgotirqsxdbsckghnihwbpyphgaixzfqyfwdsmkotdigozvagizcnxltuczpggftohkvjvackfnpsvcuhcqfmobpufrpdswferlmraokutwyqxaxdemokkvnsngtnpaggkktmxuvtaoerulvuogccnujovctucjzkjcacqayhooyjdwgbiwawkbkgernibccljapbonjaahdqhoxwajgsaitcunxsezfafcqdhmxdxfpjxzfihsvsnxueqvwtdkilmfxrvfzahfyzoqzcvujojgkljofveoizzsnqmvjwaqwxymzqfwpjyprogptaqusycqihsykfylovtwyeonvnqoldcnsqrbikgpfcobbjgzqofwvzrmhftqtwyxqbfbegpdtgbdpmndfcrtbbxpeeouafpxpvaoboctpnctewsubcbtcojqnctrbvosbpvqvyiweqyhguxdyxlfidhwxbjwevzstlxjkqjlhkrwayhixqyhcykxrdygosebppmhonzkjatmtunujbafscymfvqunyhjgrcwwexgdzhrllztstwykzkjsjsdplxbqyhfpirnfesvgvbkveuorlxttsjqtzuklqsymhbvdufevcwynymjwkbmzbeuzinscnuszzxuzilpuktqvgvnsxpbgjsvdziakykkcnnryqafpcvyafyhkyeorlcbzrnhtkjxprknobozgbypqblryuolkzoarwqfsqbtgrgefvkllrobrwbpurmjdvmngusevviadavlclvarwwtgflzsquwwgufnhoowlgwgdorbncagheifvqyfzyyovgwgbjsoqjtuvymuzjwnoktnqsrqkyjeumghiwvhrxvjkqxrhgzuiqlrpsljvuvaozqccuqyocbtuwhxypvatnuwgklwiopopynvfimftafhhsqxxwqazpupldusdgqszopdqoltdcpdrovpiieyxifqydhhqpjksbwfyldvhygexxueihfadqbfunyjgpjoohyptamxcpnlcdulgskrscopoldwekwncpottvthjqaakdpklonbzyzqszsqjfjxfdawarxtodqwtfkxxhyslrjaizzbtduncohzytwmlsianrygwsqkizsddtvwebxukuhmjbrgyewdserjkvvapknqjxqcvpzguooxtrqoeaeykiyooaqlyzuusnybchoevdwtofwrjzfcnyrajqoyfswtxiimtseshreyiuenvuwriylwqhoxwhsabacegohhjhhfihyptnllefankmployeqshrpenrdhlftvzxlxdxdgcqwnubzpglnfhymiuzufzdnzkxbuiblclpvkpmhnmgflyontqwzwhbquuuiylmmuouylewjcwyskitnefsdzxufsbpnqqujhubybpmexfzaapvrfnsstjmqwhqvmjeahsghmeejrgpeaxahrbopkvovydmabfgbbxelvikguypjyhfnhtupyjmexlaglacpxytjcdwwgtnjxpinivhsekgprdfedpiebzjtzmaoaruikqyzhgaaopmehwacdxujhpldxzumfjtrevysmgtkqibjlzspioxaspsgjlpwfduosmbssjjyasecqaydlebtwodvryjijokmxsnylrzpqozybpjissljskovpajuvdqlhrbcpqaxhiefqoxmzzkppfgbnotboizmbgitmgxzqouzedcxxplufemyqmdccrymntzfuztuodtsjquzapjsehyxidatojfeiqszlxynykgelwneechdzavwnhgtrwirfgtawycwivzulbbsrncvxizvbkzolermjxzdmrgspelucfzothhkghtqagepdsfhfwxlrtxofbgccxsijcondbbiqlelvzbgbdyytauupoeaohimzwggbnggxronfoopdzaxaeuwqtfdsqhkzkyelkjvtalmbqwsmmwjqfcswaksqdemipycfqstsxruoavetuodrtrfhucqwbaospmjnwnzctdachgvlehlmuhubqwncwscwvnwkofmbjpwmftvubrhyhruvzfpqvlfcvdikukkqxrivdwnlfimbhdnohrwzbauahwjnlsuqwljnsmfngnswypogczzgqgxcelqwgvfsomdmhhahckzfzxmjfhrdofbcjefhmisxffcoxzvtzbfnhuoverkvfkrlkqadzustxemnrkvuyagvhxirdkesmzwajbztuykqsiqknrtdhf"
	fedoraMagnet        = "magnet:?xt=urn:btih:LTLWT6I2S4R2DBLMR4YPFJBVU4LCOBMY&dn=Fedora-Workstation-Live-x86_64-38&tr=http%3A%2F%2Ftorrent.fedoraproject.org%3A6969%2Fannounce"
	complexMap          = map[string]interface{}{
		"magnetLink": fedoraMagnet,
		"count": map[string]interface{}{
			"seeders":       10053,
			"done":          2592,
			"changePerHour": -235,
		},
		"peers": []interface{}{},
		"pieces": []interface{}{
			"piece1",
			"piece2",
			3,
			map[string]interface{}{
				"sub1": "No",
				"sub2": "Yes",
			},
		},
		"downloaded": map[string]interface{}{},
	}
	complexMapTranslated = "d5:countd7:seedersi10053e4:donei2592e13:changePerHouri-235ee5:peersle6:piecesl6:piece16:piece2i3ed4:sub12:No4:sub23:Yesee10:downloadedde10:magnetLink149:magnet:?xt=urn:btih:LTLWT6I2S4R2DBLMR4YPFJBVU4LCOBMY&dn=Fedora-Workstation-Live-x86_64-38&tr=http%3A%2F%2Ftorrent.fedoraproject.org%3A6969%2Fannouncee"

	invalidTestCases = []interface{}{
		true,
		false,
		map[int]bool{},
		map[int64]string{},
		struct{ invalid string }{""},
		bytes.Buffer{},
		map[string]interface{}{
			"hello": []interface{}{
				true,
			},
		},
	}
	stringsTestCases = []string{
		"",
		":",
		"Hello World",
		"Hello:5:World",
		fedoraMagnet,
		extremelyLongString,
	}
	int64TestCases = map[int64]string{
		0:             "i0e",
		1:             "i1e",
		-1:            "i-1e",
		42316:         "i42316e",
		-2535:         "i-2535e",
		math.MaxInt64: "i9223372036854775807e",
		math.MinInt64: "i-9223372036854775808e",
	}
	intsTestCases = map[string]int{
		"i0e":                            0,
		"i1e":                            1,
		"i-1e":                           -1,
		"i42316e":                        42316,
		"i-2535e":                        -2535,
		fmt.Sprintf("i%de", math.MaxInt): math.MaxInt,
		fmt.Sprintf("i%de", math.MinInt): math.MinInt,
	}
	uintsTestCases = map[uint64]string{
		0:              "i0e",
		1:              "i1e",
		42316:          "i42316e",
		math.MaxUint64: "i18446744073709551615e",
	}
	slicesTestCases = map[string][]interface{}{
		"le":                     {},
		"llee":                   {[]interface{}{}},
		"l11:Hello Worlde":       {"Hello World"},
		"li11ele11:Hello Worlde": {11, []interface{}{}, "Hello World"},
	}
	mapTestCases = map[string]map[string]interface{}{
		"de":          {},
		"d5:hellodee": {"hello": map[string]interface{}{}},
		// "d5:hello5:world2:hii5ee": {"hello": "world", "hi": 5},
		// complexMapTranslated:      complexMap,
	}
	complexMapTestCases = map[string]map[string]interface{}{
		"de":                       {},
		"d5:hellodee":              {"hello": map[string]interface{}{}},
		"d5:hello5:world2:hii5ee":  {"hello": "world", "hi": 5},
		complexMapTranslated:       complexMap,
		"d4:name5:Alice3:agei35ee": {"name": "Alice", "age": 35},
	}
	invalidParserInputs = []string{
		"",
		"e",
		":",
		"4:abc",
		"l4:abce",
		"i:9125822395812385218357",
		"i:one",
		"li12e",
		"di12ee",
		"d",
		"1",
		"l",
		"i",
		"-1:",
		"-1:abc",
		"di3e-1:e",
		"di3e-1:abce",
	}
	invalidTypeParse = map[string]byte{
		"":       byte('x'),
		"x34l":   byte('x'),
		"1:a":    byte('s'),
		"li1ee":  byte('l'),
		"di2ee":  byte('d'),
		"i32e":   byte('i'),
		"i-123e": byte('i'),
	}
)

var (
	encoderBenchmarks = map[string]interface{}{
		"torrentString": fedoraMagnet,
		"complexMap":    complexMap,
	}
	fedoraMagnetParsed = fmt.Sprintf("%d:%s", len(fedoraMagnet), fedoraMagnet)
	parserBenchmarks   = map[string]string{
		"torrentString": fedoraMagnetParsed,
		"complexMap":    complexMapTranslated,
	}
)
