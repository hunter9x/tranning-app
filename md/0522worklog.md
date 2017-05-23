#### 午後やること:
[ ] index.html作る: mindmap list (create new)
[x] mySqlを作成
[ ] GoでDB操作するプログラムを書く
[ ] index.htmlのJsと連携させる
[ ] UIを工夫する
/*------------------------------------------------------------*/
[] 位置の同期をする
[] 削除機能をつける

#### メンタ-リング
* 明日発表練習をする
* 午後は聞かれそうなところを直す
* 理想:何のサ-ビス + `工夫`したところ + demo

## DBの設計
### 作成したマインドマップの名前一覧から開けるようにする
#### 一覧のペ-ジを作る
- DBから一覧テ-ブルの最新10件を取ってきて表示する`@一覧のペ-ジ`
- 名前が選択されたらJSONテ-ブルからそれを取ってきて`show()`に渡す
`@表示のペ-ジ`

#### マインドマップを書くために何が必要

```javascript
function load_jsmind(){
    var mind = {
        "meta":{
            "name":"demo",
            "author":"hizzgdev@163.com",
            "version":"0.2",
        },
        "format":"node_array",
        "data":[
            {"id":"root", "isroot":true, "topic":"jsMind"},

            {"id":"sub1", "parentid":"root", "topic":"sub1", "background-color":"#0000ff"},
            {"id":"sub11", "parentid":"sub1", "topic":"sub11"},
            {"id":"sub12", "parentid":"sub1", "topic":"sub12"},
            {"id":"sub13", "parentid":"sub1", "topic":"sub13"},

            {"id":"sub2", "parentid":"root", "topic":"sub2"},
            {"id":"sub21", "parentid":"sub2", "topic":"sub21"},
            {"id":"sub22", "parentid":"sub2", "topic":"sub22","foreground-color":"#33ff33"},

            {"id":"sub3", "parentid":"root", "topic":"sub3"},
        ]
    };
    var options = {
        container:'jsmind_container',
        editable:true,
        theme:'primary'
    }
    var jm = jsMind.show(options,mind);
}

load_jsmind();
```

再現するには`var jm = jsMind.show(options,mind);`
`options`では
    `container:'jsmind_container'` htmlのDOM id
    `ditable:true` 編集可能かどうか
    `theme:'primary'` themes指定

`mind`では
    `meta`宣言はいらないが`meta`の存在が必須
    `format`必須
    `data`は基本的にこういう形
        `{"id":"root", "isroot":true, "topic":"jsMind"}`root node
        `{"id":"sub3", "parentid":"root", "topic":"sub3"}`他のnode

##### 保存する(保存ボタン@js --> get_data[name, json]@js --> savemysql@go-lang)

- __デ-タはどこに保存されるか__
- __それをどうDBに書き込むかか__
    `main_data`
    * 作成した日付
    * マインドマップの名前
    * mind
    * options
    `mindmap_list`
    * 作成した日付
    * マインドマップの名前

- __どこに保存されるか__
保存形態は3つがあるらしい
* `node_tree(default)` JSON
こういう感じ⇒`options`がないじゃん?`固定すればいいかな`

```json
{"meta":{"name":"jsMind","author":"hizzgdev@163.com","version":"0.4.3"},"format":"node_tree","data":{"id":"root","topic":"jsMind Example","expanded":true}}
```

* `node_array` JSONで`data`を配列に入れる
* `freemind(.mm)` freemindに対応できるXML

- __どう保存するか__
```javascript
function save_file(){
    var mind_data = _jm.get_data();
    var mind_name = mind_data.meta.name;
    var mind_str = jsMind.util.json.json2string(mind_data);
    jsMind.util.file.save(mind_str,'text/jsmind',mind_name+'.jm');
}
```
`_jm.get_data()`でマインドマップのデ-タを保存できる(JSON)、
さらに`.meta.name`で`meta`で宣言した名前

##### DBから参照(名前選択@js --> DBから取得する@go-lang --> 表示する@js)

- __どうopenするか__
サンプルでは`choose file`してから`open file`する
`choose file`は`input type="file"`でファイルを開く?

```html
<li><input id="file_input" class="file_input" type="file"/></li>
<li><button class="sub" onclick="open_file();">open file</button></li>
```

```javascript
function open_file(){
    var file_input = document.getElementById('file_input');
    var files = file_input.files;
    if(files.length > 0){
        var file_data = files[0];
        jsMind.util.file.read(file_data,function(jsmind_data, jsmind_name){
            var mind = jsMind.util.json.string2json(jsmind_data);
            if(!!mind){
                _jm.show(mind);
            }else{
                prompt_info('can not open this file as mindmap');
            }
        });
    }else{
        prompt_info('please choose a file first')
    }
}
```

`mind`(JSON)を使えばいい
その前に`jsMind.util.json.string2json`でstringに変換する
⇒DBから持ってくる時は`string`にする必要ある


### DBサ-バ-を作るかどうか
先に`local DB`で作成するべき
#### mySql installing:
```bash
# pwdはdocker userと同じものにした
apt-get -y install mysql-server-5.7

# user list:
mysql> select user,host from mysql.user;
+------------------+-----------+
| user             | host      |
+------------------+-----------+
| debian-sys-maint | localhost |
| mysql.sys        | localhost |
| root             | localhost |
+------------------+-----------+
```

#### DB creating:
##### DB名: `mindmap`
```sql
CREATE DATABASE mindmap;
```
##### テ-ブル作成

`main_data`
* id
* 作成した日付
* マインドマップの名前
* mind(JSON保存)
* options

```sql
CREATE TABLE `mind_data` (
        `mapID` int NOT NULL,
        `created` DATETIME NOT NULL,
        `name` varchar(20) NOT NULL,
        `mind` varchar(1024),
        PRIMARY KEY (`mapID`)
        );
```
JSONをどうやってmySqlに保存するか?`varchar`で保存するか

`mindmap_list`**要らない**
* id
* マインドマップの名前

#### Bootstrapを使っていい感じのindex.htmlを作ろう
`Bootstrap CDN `の方法を選択する
```html
<!DOCTYPE html>
<html lang="ja">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>タイトル</title>
    <link href="css/bootstrap.min.css" rel="stylesheet">
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
  </head>
  <body>
    <h1>Hello, world!</h1>
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js"></script>
    <script src="js/bootstrap.min.js"></script>
  </body>
</html>
```

##### bootstrapの使い方を勉強する
`コンテナ`はペ-ジの基本シ-トとなる要素
Bootstrap の各要素をコンテナの中に記述
`container` は固定的コンテナを生成します
ウィンドウの横幅に応じて段階的に横幅が変動します
とにかく、containerの中で全部書けばいいかな

[ICON photo](https://pbs.twimg.com/profile_images/642108475980574720/TL3Bgv-C.jpg)
__上から各ペ-ジを作って行く__
[x]`index.html`(create new & view past)
    + ボタンにクリックする時に参照するペ-ジはまだ指定していない
    + ボタンのデザインをもうちょっと簡潔にしたい
    [] `html`ペ-ジをデザインする
[ ]`view.html`(DBから10個持ってくる)
    [] `html` mindmapの一覧リスト
    [] `javascript`
        * ペ-ジロ-ドされたら`list`を取る命令をgoに送る
        * 選択された`name`をgo-langに渡す
    [] `go-lang` DBから
        * 一覧取得命令をもらい、DBから`fetch`する
        * `name`の情報を持ってくる

[ ]`create new`(作成画面に`保存`とデザイン改良)
    [] `html` デザインの改良
    [] `javascript` 保存ボタンを追加、goと連携
    [] `go-lang` jsから値をもらってmysqlに保存する

### これってAPIサ-バ-かな
#####`API`と`UI`の区別:
* `人`が触るのがUI(ボタンなど)
* APIは`プログラム`が操作できるもの

##### APIとUIが連携した仕組み
[UI](/mindmaps/1)
[API](mindmap1を返す)
UIから`Request`をAPIに送る→デ-タが返ってくる(JSON等)
__ほしいものを文字のrequestでgetできる__

#### これをAPIにするためには
1. 複数の`session`と`html/js`の連携?