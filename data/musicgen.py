import json
import pickle
import re

from fuzzywuzzy import process
from jikanpy import Jikan


class Aninx:
    Endpoint = "https://aninx.com"
    Data = []


def remove_prefix(text, prefix):
    return text[text.startswith(prefix) and len(prefix):].rstrip("\n")


def get_top_500():
    data = set()
    jikan = Jikan()

    for i in range(1, 11):
        print(f"Processing page: {i}")
        if i == 5:
            continue
        top_anime = jikan.top(type='anime', page=i, subtype='bypopularity')
        print(len(top_anime["top"]))
        for j in range(0, 50):
            if top_anime["top"][j]["type"] not in ["OVA", "Music", "Special"]:
                data.add(top_anime["top"][j]["title"])

    data5 = '{"request_hash":"request:top:abde8908f934dc2d681a07b8cafca26c8e9352ce","request_cached":true,"request_cache_expiry":24705,"top":[{"mal_id":31798,"rank":201,"title":"Kiznaiver","url":"https:\/\/myanimelist.net\/anime\/31798\/Kiznaiver","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/6\/78466.jpg?s=9d7ab9c0d0680bc6f5717f722d78eb96","type":"TV","episodes":12,"start_date":"Apr 2016","end_date":"Jun 2016","members":432784,"score":7.5},{"mal_id":36474,"rank":202,"title":"Sword Art Online: Alicization","url":"https:\/\/myanimelist.net\/anime\/36474\/Sword_Art_Online__Alicization","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/1993\/93837.jpg?s=5054e1f2d07192aad8470d3992640f6d","type":"TV","episodes":24,"start_date":"Oct 2018","end_date":"Mar 2019","members":430808,"score":7.64},{"mal_id":22789,"rank":203,"title":"Barakamon","url":"https:\/\/myanimelist.net\/anime\/22789\/Barakamon","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/12\/65427.jpg?s=fa2f22ba147730525acb265769e568ab","type":"TV","episodes":12,"start_date":"Jul 2014","end_date":"Sep 2014","members":426809,"score":8.44},{"mal_id":23281,"rank":204,"title":"Psycho-Pass 2","url":"https:\/\/myanimelist.net\/anime\/23281\/Psycho-Pass_2","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/1197\/100616.jpg?s=33bb43f9b833bfb920260791254cffd7","type":"TV","episodes":11,"start_date":"Oct 2014","end_date":"Dec 2014","members":424381,"score":7.47},{"mal_id":9041,"rank":205,"title":"IS: Infinite Stratos","url":"https:\/\/myanimelist.net\/anime\/9041\/IS__Infinite_Stratos","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/3\/74045.jpg?s=858bc6fb1800d917813882d845ead68d","type":"TV","episodes":12,"start_date":"Jan 2011","end_date":"Apr 2011","members":424377,"score":6.89},{"mal_id":30307,"rank":206,"title":"Monster Musume no Iru Nichijou","url":"https:\/\/myanimelist.net\/anime\/30307\/Monster_Musume_no_Iru_Nichijou","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/9\/75104.jpg?s=c762b931d2781cc766c92a9eb62cb5b1","type":"TV","episodes":12,"start_date":"Jul 2015","end_date":"Sep 2015","members":423837,"score":7.18},{"mal_id":32,"rank":207,"title":"Neon Genesis Evangelion: The End of Evangelion","url":"https:\/\/myanimelist.net\/anime\/32\/Neon_Genesis_Evangelion__The_End_of_Evangelion","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/12\/39305.jpg?s=28e24089c08d7b088d7cb1c259a4cbf3","type":"Movie","episodes":1,"start_date":"Jul 1997","end_date":"Jul 1997","members":422914,"score":8.51},{"mal_id":14075,"rank":208,"title":"Zetsuen no Tempest","url":"https:\/\/myanimelist.net\/anime\/14075\/Zetsuen_no_Tempest","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/7\/42453.jpg?s=f4722250dda9e87791ced4daec787781","type":"TV","episodes":24,"start_date":"Oct 2012","end_date":"Mar 2013","members":422701,"score":8.05},{"mal_id":32542,"rank":209,"title":"Sakamoto Desu ga?","url":"https:\/\/myanimelist.net\/anime\/32542\/Sakamoto_Desu_ga","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/4\/79468.jpg?s=7373e097b2fffef594362c7e3466b065","type":"TV","episodes":12,"start_date":"Apr 2016","end_date":"Jul 2016","members":421690,"score":7.64},{"mal_id":355,"rank":210,"title":"Shakugan no Shana","url":"https:\/\/myanimelist.net\/anime\/355\/Shakugan_no_Shana","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/8\/21197.jpg?s=17405ca800c1c6fe7a9da573c74f6a79","type":"TV","episodes":24,"start_date":"Oct 2005","end_date":"Mar 2006","members":420615,"score":7.57},{"mal_id":7674,"rank":211,"title":"Bakuman.","url":"https:\/\/myanimelist.net\/anime\/7674\/Bakuman","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/6\/26138.jpg?s=d295be0db62defdde7588f65acce433a","type":"TV","episodes":25,"start_date":"Oct 2010","end_date":"Apr 2011","members":419880,"score":8.27},{"mal_id":18277,"rank":212,"title":"Strike the Blood","url":"https:\/\/myanimelist.net\/anime\/18277\/Strike_the_Blood","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/5\/56163.jpg?s=d5ceae087d157ee8af8fed5d1fede8de","type":"TV","episodes":24,"start_date":"Oct 2013","end_date":"Mar 2014","members":419707,"score":7.24},{"mal_id":28927,"rank":213,"title":"Owari no Seraph: Nagoya Kessen-hen","url":"https:\/\/myanimelist.net\/anime\/28927\/Owari_no_Seraph__Nagoya_Kessen-hen","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/9\/76632.jpg?s=bae15204f2080371d7ce0e712a245ecd","type":"TV","episodes":12,"start_date":"Oct 2015","end_date":"Dec 2015","members":419237,"score":7.7},{"mal_id":790,"rank":214,"title":"Ergo Proxy","url":"https:\/\/myanimelist.net\/anime\/790\/Ergo_Proxy","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/11\/6259.jpg?s=4fd55e1e59f0bb02ea88a28ed88834c1","type":"TV","episodes":23,"start_date":"Feb 2006","end_date":"Aug 2006","members":417490,"score":7.94},{"mal_id":33487,"rank":215,"title":"Masamune-kun no Revenge","url":"https:\/\/myanimelist.net\/anime\/33487\/Masamune-kun_no_Revenge","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/12\/83709.jpg?s=27cac4795a28ba1b49cdd20fb1528841","type":"TV","episodes":12,"start_date":"Jan 2017","end_date":"Mar 2017","members":417234,"score":6.93},{"mal_id":28497,"rank":216,"title":"Rokka no Yuusha","url":"https:\/\/myanimelist.net\/anime\/28497\/Rokka_no_Yuusha","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/9\/74374.jpg?s=d6dbfb61c51a8322468af35a8deee207","type":"TV","episodes":12,"start_date":"Jul 2015","end_date":"Sep 2015","members":416801,"score":7.43},{"mal_id":7311,"rank":217,"title":"Suzumiya Haruhi no Shoushitsu","url":"https:\/\/myanimelist.net\/anime\/7311\/Suzumiya_Haruhi_no_Shoushitsu","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/2\/73842.jpg?s=b203b8eab5c8f6752221df27efa7fbf8","type":"Movie","episodes":1,"start_date":"Feb 2010","end_date":"Feb 2010","members":414030,"score":8.69},{"mal_id":17074,"rank":218,"title":"Monogatari Series: Second Season","url":"https:\/\/myanimelist.net\/anime\/17074\/Monogatari_Series__Second_Season","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/3\/52133.jpg?s=b22811f98b26399cd65787af1447848c","type":"TV","episodes":26,"start_date":"Jul 2013","end_date":"Dec 2013","members":413411,"score":8.78},{"mal_id":3457,"rank":219,"title":"Vampire Knight","url":"https:\/\/myanimelist.net\/anime\/3457\/Vampire_Knight","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/3\/7327.jpg?s=9845f0265ce9f330bb2f4e45a30cdc50","type":"TV","episodes":13,"start_date":"Apr 2008","end_date":"Jul 2008","members":412908,"score":7.15},{"mal_id":22147,"rank":220,"title":"Amagi Brilliant Park","url":"https:\/\/myanimelist.net\/anime\/22147\/Amagi_Brilliant_Park","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/5\/85435.jpg?s=7ed74447820deab77b52178da0168a5d","type":"TV","episodes":13,"start_date":"Oct 2014","end_date":"Dec 2014","members":412808,"score":7.58},{"mal_id":43,"rank":221,"title":"Koukaku Kidoutai","url":"https:\/\/myanimelist.net\/anime\/43\/Koukaku_Kidoutai","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/10\/82594.jpg?s=077dce74e000ea7dcd77c738bd1cbee0","type":"Movie","episodes":1,"start_date":"Nov 1995","end_date":"Nov 1995","members":412502,"score":8.31},{"mal_id":7791,"rank":222,"title":"K-On!!","url":"https:\/\/myanimelist.net\/anime\/7791\/K-On","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/12\/76121.jpg?s=0b734a44ed13694c2396e1eb8260e0bf","type":"TV","episodes":26,"start_date":"Apr 2010","end_date":"Sep 2010","members":411590,"score":8.15},{"mal_id":32995,"rank":223,"title":"Yuri!!! on Ice","url":"https:\/\/myanimelist.net\/anime\/32995\/Yuri_on_Ice","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/6\/81149.jpg?s=4fe09eba4eb9a04c30a69f523778551d","type":"TV","episodes":12,"start_date":"Oct 2016","end_date":"Dec 2016","members":408742,"score":8.02},{"mal_id":5630,"rank":224,"title":"Higashi no Eden","url":"https:\/\/myanimelist.net\/anime\/5630\/Higashi_no_Eden","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/9\/15033.jpg?s=4a5a58720c1b20212230343e9297b9df","type":"TV","episodes":11,"start_date":"Apr 2009","end_date":"Jun 2009","members":407134,"score":7.88},{"mal_id":3455,"rank":225,"title":"To LOVE-Ru","url":"https:\/\/myanimelist.net\/anime\/3455\/To_LOVE-Ru","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/13\/22544.jpg?s=c433699f1bb4e7243754703dff0bdb32","type":"TV","episodes":26,"start_date":"Apr 2008","end_date":"Sep 2008","members":406958,"score":7.14},{"mal_id":23321,"rank":226,"title":"Log Horizon 2nd Season","url":"https:\/\/myanimelist.net\/anime\/23321\/Log_Horizon_2nd_Season","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/5\/68097.jpg?s=843e25fefdf0cffa49ba7bfb41e2d8d3","type":"TV","episodes":25,"start_date":"Oct 2014","end_date":"Mar 2015","members":404540,"score":7.64},{"mal_id":7724,"rank":227,"title":"Shiki","url":"https:\/\/myanimelist.net\/anime\/7724\/Shiki","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/8\/21343.jpg?s=a46bd1834acc06a62522a68ce603bf22","type":"TV","episodes":22,"start_date":"Jul 2010","end_date":"Dec 2010","members":403536,"score":7.84},{"mal_id":1604,"rank":228,"title":"Katekyo Hitman Reborn!","url":"https:\/\/myanimelist.net\/anime\/1604\/Katekyo_Hitman_Reborn","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/13\/80000.jpg?s=b52850d927668c7331d1d0807e20a886","type":"TV","episodes":203,"start_date":"Oct 2006","end_date":"Sep 2010","members":402974,"score":8.22},{"mal_id":18671,"rank":229,"title":"Chuunibyou demo Koi ga Shitai! Ren","url":"https:\/\/myanimelist.net\/anime\/18671\/Chuunibyou_demo_Koi_ga_Shitai_Ren","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/7\/56643.jpg?s=8802dedde3f74fa975ec9e7b7ce4baed","type":"TV","episodes":12,"start_date":"Jan 2014","end_date":"Mar 2014","members":402868,"score":7.57},{"mal_id":26055,"rank":230,"title":"JoJo no Kimyou na Bouken Part 3: Stardust Crusaders 2nd Season","url":"https:\/\/myanimelist.net\/anime\/26055\/JoJo_no_Kimyou_na_Bouken_Part_3__Stardust_Crusaders_2nd_Season","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/11\/75045.jpg?s=f50d24a2193eebe6b4cf0324cc93333c","type":"TV","episodes":24,"start_date":"Jan 2015","end_date":"Jun 2015","members":402006,"score":8.54},{"mal_id":24415,"rank":231,"title":"Kuroko no Basket 3rd Season","url":"https:\/\/myanimelist.net\/anime\/24415\/Kuroko_no_Basket_3rd_Season","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/4\/68299.jpg?s=047c586d552e0b0ea90cb19074fb6ea8","type":"TV","episodes":25,"start_date":"Jan 2015","end_date":"Jun 2015","members":401298,"score":8.44},{"mal_id":32729,"rank":232,"title":"Orange","url":"https:\/\/myanimelist.net\/anime\/32729\/Orange","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/4\/80110.jpg?s=4400af9791f67ab556580a6d984c0233","type":"TV","episodes":13,"start_date":"Jul 2016","end_date":"Sep 2016","members":399268,"score":7.66},{"mal_id":15315,"rank":233,"title":"Mondaiji-tachi ga Isekai kara Kuru Sou Desu yo?","url":"https:\/\/myanimelist.net\/anime\/15315\/Mondaiji-tachi_ga_Isekai_kara_Kuru_Sou_Desu_yo","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/12\/43369.jpg?s=f75c89ce4dcb7e8ad77a7d55124a2a97","type":"TV","episodes":10,"start_date":"Jan 2013","end_date":"Mar 2013","members":398343,"score":7.66},{"mal_id":30694,"rank":234,"title":"Dragon Ball Super","url":"https:\/\/myanimelist.net\/anime\/30694\/Dragon_Ball_Super","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/7\/74606.jpg?s=45a02b9bcc08d906381a086129f6a0f5","type":"TV","episodes":131,"start_date":"Jul 2015","end_date":"Mar 2018","members":397282,"score":7.55},{"mal_id":35120,"rank":235,"title":"Devilman: Crybaby","url":"https:\/\/myanimelist.net\/anime\/35120\/Devilman__Crybaby","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/2\/89973.jpg?s=717c05245eb3166f30bac5945b9cc6c7","type":"ONA","episodes":10,"start_date":"Jan 2018","end_date":"Jan 2018","members":396844,"score":7.91},{"mal_id":31933,"rank":236,"title":"JoJo no Kimyou na Bouken Part 4: Diamond wa Kudakenai","url":"https:\/\/myanimelist.net\/anime\/31933\/JoJo_no_Kimyou_na_Bouken_Part_4__Diamond_wa_Kudakenai","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/3\/79156.jpg?s=9189ddeb4a702d774b985fe16f315a81","type":"TV","episodes":39,"start_date":"Apr 2016","end_date":"Dec 2016","members":395328,"score":8.58},{"mal_id":33506,"rank":237,"title":"Ao no Exorcist: Kyoto Fujouou-hen","url":"https:\/\/myanimelist.net\/anime\/33506\/Ao_no_Exorcist__Kyoto_Fujouou-hen","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/5\/85201.jpg?s=98f1a88ffabf7c0295847361307d0857","type":"TV","episodes":12,"start_date":"Jan 2017","end_date":"Mar 2017","members":395122,"score":7.55},{"mal_id":16067,"rank":238,"title":"Nagi no Asu kara","url":"https:\/\/myanimelist.net\/anime\/16067\/Nagi_no_Asu_kara","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/7\/53549.jpg?s=60d473da16a7ba46e66caf1c1f796cbd","type":"TV","episodes":26,"start_date":"Oct 2013","end_date":"Apr 2014","members":394796,"score":8.17},{"mal_id":14289,"rank":239,"title":"Sukitte Ii na yo.","url":"https:\/\/myanimelist.net\/anime\/14289\/Sukitte_Ii_na_yo","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/11\/39777.jpg?s=2e1e9b772105988e635fd0f8a18c8248","type":"TV","episodes":13,"start_date":"Oct 2012","end_date":"Dec 2012","members":392513,"score":7.57},{"mal_id":225,"rank":240,"title":"Dragon Ball GT","url":"https:\/\/myanimelist.net\/anime\/225\/Dragon_Ball_GT","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/11\/20727.jpg?s=f142e2656af64230f12c64c96b5af3fa","type":"TV","episodes":64,"start_date":"Feb 1996","end_date":"Nov 1997","members":391469,"score":6.56},{"mal_id":11499,"rank":241,"title":"Sankarea","url":"https:\/\/myanimelist.net\/anime\/11499\/Sankarea","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/1487\/95651.jpg?s=45618dace46cff103410af6be9c07823","type":"TV","episodes":12,"start_date":"Apr 2012","end_date":"Jun 2012","members":391045,"score":7.42},{"mal_id":32901,"rank":242,"title":"Eromanga-sensei","url":"https:\/\/myanimelist.net\/anime\/32901\/Eromanga-sensei","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/2\/86468.jpg?s=1b0ab8cd10df83ba9f1de708d80179c6","type":"TV","episodes":12,"start_date":"Apr 2017","end_date":"Jun 2017","members":389649,"score":6.68},{"mal_id":22729,"rank":243,"title":"Aldnoah.Zero","url":"https:\/\/myanimelist.net\/anime\/22729\/AldnoahZero","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/7\/60263.jpg?s=24207466eff27c10e008061f265b2bca","type":"TV","episodes":12,"start_date":"Jul 2014","end_date":"Sep 2014","members":388866,"score":7.49},{"mal_id":9513,"rank":244,"title":"Beelzebub","url":"https:\/\/myanimelist.net\/anime\/9513\/Beelzebub","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/3\/28013.jpg?s=33c1e4b2fc322189b6799b171e342a04","type":"TV","episodes":60,"start_date":"Jan 2011","end_date":"Mar 2012","members":388722,"score":7.97},{"mal_id":28297,"rank":245,"title":"Ore Monogatari!!","url":"https:\/\/myanimelist.net\/anime\/28297\/Ore_Monogatari","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/13\/69455.jpg?s=25a03f3acb40aeebb019fe11123b38cd","type":"TV","episodes":24,"start_date":"Apr 2015","end_date":"Sep 2015","members":388536,"score":8.02},{"mal_id":27787,"rank":246,"title":"Nisekoi:","url":"https:\/\/myanimelist.net\/anime\/27787\/Nisekoi_","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/13\/72626.jpg?s=4bb62a600ef1886a6b9d3bef8ba967e6","type":"TV","episodes":12,"start_date":"Apr 2015","end_date":"Jun 2015","members":387242,"score":7.44},{"mal_id":31765,"rank":247,"title":"Sword Art Online Movie: Ordinal Scale","url":"https:\/\/myanimelist.net\/anime\/31765\/Sword_Art_Online_Movie__Ordinal_Scale","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/4\/83811.jpg?s=92313c5ef2e6df1fd837384c9539c428","type":"Movie","episodes":1,"start_date":"Feb 2017","end_date":"Feb 2017","members":384451,"score":7.82},{"mal_id":339,"rank":248,"title":"Serial Experiments Lain","url":"https:\/\/myanimelist.net\/anime\/339\/Serial_Experiments_Lain","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/1718\/91550.jpg?s=a7345c87c4b2e5632a0562f3204d52bb","type":"TV","episodes":13,"start_date":"Jul 1998","end_date":"Sep 1998","members":384235,"score":8.01},{"mal_id":37675,"rank":249,"title":"Overlord III","url":"https:\/\/myanimelist.net\/anime\/37675\/Overlord_III","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/1511\/93473.jpg?s=c1ff35d315aa8ee8ba1b23a1307c2af4","type":"TV","episodes":13,"start_date":"Jul 2018","end_date":"Oct 2018","members":383806,"score":8.06},{"mal_id":30544,"rank":250,"title":"Gakusen Toshi Asterisk","url":"https:\/\/myanimelist.net\/anime\/30544\/Gakusen_Toshi_Asterisk","image_url":"https:\/\/cdn.myanimelist.net\/images\/anime\/5\/76034.jpg?s=f6512ca583dffe851392ded9316fef15","type":"TV","episodes":12,"start_date":"Oct 2015","end_date":"Dec 2015","members":382794,"score":7.01}]}'
    parsed_json = json.loads(data5)

    for j in range(0, 50):
        if parsed_json["top"][j]["type"] not in ["OVA", "Music", "Special"]:
            print(parsed_json["top"][j]["title"])
            data.add(parsed_json["top"][j]["title"])

    for i in range(11, 16):
        print(f"Processing page: {i}")
        top_anime = jikan.top(type='anime', page=i, subtype='bypopularity')
        for j in range(0, 50):
            if (top_anime["top"][j]["type"] not in ["OVA", "Music", "Special"]) and top_anime["top"][j]["score"] >= 7.5:
                data.add(top_anime["top"][j]["title"])

    for i in range(16, 21):
        print(f"Processing page: {i}")
        top_anime = jikan.top(type='anime', page=i, subtype='bypopularity')
        for j in range(0, 50):
            if (top_anime["top"][j]["type"] not in ["OVA", "Music", "Special"]) and top_anime["top"][j]["score"] >= 8:
                data.add(top_anime["top"][j]["title"])

    return data


if __name__ == "__main__":
    folders = list(range(2000, 2021)) + ["60s", "70s", "80s", "90s", "misc"]
    #
    # for i in folders:
    #     r = requests.get(f"{Aninx.Endpoint}/{i}/success.txt")
    #     with open(f"{i}success.txt", 'wb') as f:
    #         f.write(r.content)
    #
    for i in folders:
        with open(f"{i}success.txt", "r", encoding="utf8", errors="ignore") as f:
            line = f.readline()
            while line:
                search = re.search("^.{7}─ (.+)\n$", line)
                if search:
                    Aninx.Data.append({"name": search[1], "songs": []})
                else:
                    search = re.search("^.{10}─ (.+)\n$", line)
                    if search:
                        Aninx.Data[-1]["songs"].append({"songname": search[1]})
                    else:
                        search = re.search("^.{12}└─ 0: (.+)\n$", line)
                        if search:
                            Aninx.Data[-1]["songs"][-1]["url"] = search[1]
                line = f.readline()

    # filter_set = get_top_500()
    # print(len(filter_set))
    # with open('filter_set.pickle', 'wb') as handle:
    #     pickle.dump(filter_set, handle, protocol=pickle.HIGHEST_PROTOCOL)

    with open('filter_set.pickle', 'rb') as handle:
        filter_set = pickle.load(handle)

    new_data = []
    count = 0
    print(f"Total: {len(Aninx.Data)}")
    for anime in Aninx.Data:
        count += 1
        result = process.extract(anime['name'], filter_set, limit=1)[0]
        if result[1] == 100:
            new_data.append(anime)
        else:
            if result[0].endswith(" (TV)") and result[1] >= 90:
                new_data.append(anime)
        if count % 200 == 0:
            print(count)

    print(len(new_data))

    with open(f"dataset_filtered.json", "w") as f:
        json.dump(new_data, f)
