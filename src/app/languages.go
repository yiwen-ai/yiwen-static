package app

import (
	"encoding/json"
	"strings"
)

var Languages [][]string

func init() {
	if err := json.Unmarshal([]byte(languages), &Languages); err != nil {
		panic(err)
	}
}

func Lang639_3(lang string) string {
	lang = strings.ToLower(lang)

	for _, vv := range Languages {
		if vv[0] == lang || vv[1] == lang || strings.ToLower(vv[2]) == lang || strings.ToLower(vv[3]) == lang {
			return vv[1]
		}
	}

	return ""
}

var languages = `[["ab","abk","Abkhazian","аҧсуа бызшәа\u200e (Aṗsua byzšwa), аҧсшәа\u200e (Aṗsšwa)"],["aa","aar","Afar","Qafar"],["af","afr","Afrikaans","Afrikaans"],["ak","aka","Akan","Akan"],["sq","sqi","Albanian","shqip"],["am","amh","Amharic","አማርኛ"],["ar","ara","Arabic","العربية"],["an","arg","Aragonese","Aragonés"],["hy","hye","Armenian","հայերեն"],["as","asm","Assamese","অসমীয়া"],["av","ava","Avaric","авар мацӏ\u200e (Awar mac̣), магӏарул мацӏ\u200e (Maʿarul mac̣)"],["az","aze","Azerbaijani","azərbaycan"],["bm","bam","Bambara","bamanakan"],["ba","bak","Bashkir","башҡорт теле\u200e (Başķort tele), башҡортса\u200e (Başķortsa)"],["eu","eus","Basque","euskara"],["be","bel","Belarusian","беларуская"],["bn","ben","Bengali","বাংলা"],["bi","bis","Bislama","Bislama"],["bs","bos","Bosnian","bosanski"],["br","bre","Breton","brezhoneg"],["bg","bul","Bulgarian","български"],["my","mya","Burmese","မြန်မာ"],["ca","cat","Catalan","català"],["ch","cha","Chamorro","Chamorru"],["ce","che","Chechen","нохчийн"],["zh","zho","Chinese","中文"],["cu","chu","Church Slavic","церковнослове́нскїй"],["cv","chv","Chuvash","Чӑвашла"],["kw","cor","Cornish","kernewek"],["co","cos","Corsican","Corsu"],["hr","hrv","Croatian","hrvatski"],["cs","ces","Czech","čeština"],["da","dan","Danish","dansk"],["dv","div","Dhivehi","ދިވެހިބަސް"],["nl","nld","Dutch","Nederlands"],["dz","dzo","Dzongkha","རྫོང་ཁ"],["en","eng","English","English"],["eo","epo","Esperanto","esperanto"],["et","est","Estonian","eesti"],["ee","ewe","Ewe","Eʋegbe"],["fo","fao","Faroese","føroyskt"],["fi","fin","Finnish","suomi"],["fr","fra","French","français"],["ff","ful","Fulah","Pulaar"],["gl","glg","Galician","galego"],["lg","lug","Ganda","Luganda"],["ka","kat","Georgian","ქართული"],["de","deu","German","Deutsch"],["gu","guj","Gujarati","ગુજરાતી"],["ht","hat","Haitian","Ayisyen, Kreyòl"],["ha","hau","Hausa","Hausa"],["he","heb","Hebrew","עברית"],["hi","hin","Hindi","हिन्दी"],["hu","hun","Hungarian","magyar"],["is","isl","Icelandic","íslenska"],["ig","ibo","Igbo","Igbo"],["id","ind","Indonesian","Indonesia"],["ia","ina","Interlingua","interlingua"],["iu","iku","Inuktitut","ᐃᓄᒃᑎᑐᑦ"],["ga","gle","Irish","Gaeilge"],["it","ita","Italian","italiano"],["ja","jpn","Japanese","日本語"],["jv","jav","Javanese","Basa Jawa"],["kl","kal","Kalaallisut","kalaallisut"],["kn","kan","Kannada","ಕನ್ನಡ"],["ks","kas","Kashmiri","کٲشُر"],["kk","kaz","Kazakh","қазақ тілі"],["km","khm","Khmer","ខ្មែរ"],["ki","kik","Kikuyu","Gikuyu"],["rw","kin","Kinyarwanda","Kinyarwanda"],["ky","kir","Kirghiz","кыргызча"],["ko","kor","Korean","한국어"],["kj","kua","Kuanyama","Oshikwanyama"],["ku","kur","Kurdish","kurdî"],["lo","lao","Lao","ລາວ"],["lv","lav","Latvian","latviešu"],["li","lim","Limburgan","Limbourgeois, Limburgs"],["ln","lin","Lingala","lingála"],["lt","lit","Lithuanian","lietuvių"],["lu","lub","Luba-Katanga","Tshiluba"],["lb","ltz","Luxembourgish","Lëtzebuergesch"],["mk","mkd","Macedonian","македонски"],["mg","mlg","Malagasy","Malagasy"],["ms","msa","Malay","Melayu"],["ml","mal","Malayalam","മലയാളം"],["mt","mlt","Maltese","Malti"],["gv","glv","Manx","Gaelg"],["mi","mri","Maori","te reo Maori"],["mr","mar","Marathi","मराठी"],["el","ell","Modern Greek","Ελληνικά"],["mn","mon","Mongolian","монгол"],["nv","nav","Navajo","Diné"],["ne","nep","Nepali","नेपाली"],["nd","nde","North Ndebele","isiNdebele"],["se","sme","Northern Sami","davvisámegiella"],["no","nor","Norwegian","Norsk"],["nn","nno","Norwegian Nynorsk","nynorsk"],["ny","nya","Nyanja","Chichewa, chiCheŵa\u200e (Chichewa), chiNyanja"],["oc","oci","Occitan","occitan"],["or","ori","Oriya","ଓଡ଼ିଆ"],["om","orm","Oromo","Oromoo"],["os","oss","Ossetian","ирон"],["pa","pan","Panjabi","ਪੰਜਾਬੀ"],["fa","fas","Persian","فارسی"],["pl","pol","Polish","polski"],["pt","por","Portuguese","português"],["ps","pus","Pushto","پښتو"],["qu","que","Quechua","Runasimi"],["ro","ron","Romanian","română"],["rm","roh","Romansh","rumantsch"],["rn","run","Rundi","Ikirundi"],["ru","rus","Russian","русский"],["sm","smo","Samoan","Gagana Samoa"],["sg","sag","Sango","Sängö"],["sa","san","Sanskrit","संस्कृत भाषा"],["gd","gla","Scottish Gaelic","Gàidhlig"],["sr","srp","Serbian","српски"],["sn","sna","Shona","chiShona"],["ii","iii","Sichuan Yi","ꆈꌠꉙ"],["sd","snd","Sindhi","سنڌي"],["si","sin","Sinhala","සිංහල"],["sk","slk","Slovak","slovenčina"],["sl","slv","Slovenian","slovenščina"],["so","som","Somali","Soomaali"],["nr","nbl","South Ndebele","isiNdebele"],["st","sot","Southern Sotho","Sesotho"],["es","spa","Spanish","español"],["su","sun","Sundanese","Sunda"],["sw","swa","Swahili","Kiswahili"],["ss","ssw","Swati","siSwati"],["sv","swe","Swedish","svenska"],["tl","tgl","Tagalog","Tagalog"],["ty","tah","Tahitian","Reo Tahiti"],["tg","tgk","Tajik","тоҷикӣ"],["ta","tam","Tamil","தமிழ்"],["tt","tat","Tatar","татар"],["te","tel","Telugu","తెలుగు"],["th","tha","Thai","ไทย"],["bo","bod","Tibetan","བོད་སྐད་"],["ti","tir","Tigrinya","ትግርኛ"],["to","ton","Tonga","lea fakatonga"],["ts","tso","Tsonga","Xitsonga"],["tn","tsn","Tswana","Setswana"],["tr","tur","Turkish","Türkçe"],["tk","tuk","Turkmen","Türkmen dili"],["ug","uig","Uighur","ئۇيغۇرچە"],["uk","ukr","Ukrainian","українська"],["ur","urd","Urdu","اردو"],["uz","uzb","Uzbek","o‘zbek"],["ve","ven","Venda","Tshivenḓa"],["vi","vie","Vietnamese","Tiếng Việt"],["cy","cym","Welsh","Cymraeg"],["fy","fry","Western Frisian","Frysk"],["wo","wol","Wolof","Wolof"],["xh","xho","Xhosa","isiXhosa"],["yi","yid","Yiddish","ייִדיש"],["yo","yor","Yoruba","Èdè Yorùbá"],["zu","zul","Zulu","isiZulu"]]`
