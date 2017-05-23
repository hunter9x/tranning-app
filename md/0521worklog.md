## Node location sync
### No parent & no id err
#### lib formatと違うか?
そんなことはなさそう.ただ，`isroot: true` と `id:root` となっている(default)
-> `node_id`を引数にしない時と同じじゃん
#### じゃあ，引数にした時と何か状況が変わっているかな
#### jQuery: ready & loadの実行順番
`ready`は以下の関数は全て同じタイミング，つまり **HTMLがloadされた直後**
```javascript
$(function(){
  //処理
});
jQuery(function(){
  //処理
});
jQuery(document).ready(function(){
  //処理
});
```
`load`はWindow内の全て(画像等も含めて)の読み込みが完了した時に呼び出され，要は **一番最後**
```javascript
$(window).load(function() {});
```
#### Websocketが悪い
##### Libの仕様:
`msg`でnodeを作ってから`parent`などが作られる -> `どこに追加するからは`*画面操作*`で指定` デ-タが生成された訳ではない

```javascript
        click_handle:function(e){
            if (!this.options.default_event_handle['enable_click_handle']) {
                return;
            }
            var element = e.target || event.srcElement;
            var isexpander = this.view.is_expander(element);
            if(isexpander){
                var nodeid = this.view.get_binded_nodeid(element);
                if(!!nodeid){
                    this.toggle_node(nodeid);
                }
            }
        },
```

##### Websocketを使ったせいで
`msg`と生成されていない`node-id`を *同時に* 送ってしまう -> `undefine node`