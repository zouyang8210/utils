
    let pageNo = 0;
    let loading = true;
    let  url="https://www.goshop.ink:8080/";
    let app = new Vue({
        el: '#app',
        data: {
            sites: [],
            keyName:"",
            cmd:"",
            hint:"复制淘口令",
            tabBarImgArr: [   //图片切换
                {normal: './static/img/buy.png', selected: './static/img/buy_s.png'},
                {normal: './static/img/coupon.png', selected: './static/img/coupon_s.png'},
            ],
            commShow:true,
            buyShow:false,
            topShow:false,
            cid:"",

        },
       watch:{
           sites:function () {//sites数据渲染完成
             this.$nextTick(function () {
                 loading=true;
             });
           },
           keyName:function () {
               pageNo=0;
           }

       },
        methods: {
            //类目ID获取信息
            getCatCommodidy:function (cat) {
                pageNo++;
                this.$http.get(url+'cat?cat='+this.cid+'&pageNo=' + pageNo).then(function (res) {
                        if(pageNo==1){
                            this.sites = res.body.result_list.map_data;
                        }else{
                            this.sites = this.sites.concat(res.body.result_list.map_data);
                        }
                    }
                );
            },
            //默认页信息
            defaultPage: function () {
                pageNo++;
                this.$http.get(url+'default?pageNo=' + pageNo).then(function (res) {
                        // console.log(res.body);
                        if (pageNo == 1) {
                            this.sites = res.body.result_list.map_data;
                        } else {
                            this.sites = this.sites.concat(res.body.result_list.map_data);
                        }
                    }
                );
            },
            //搜索商品信息
            search: function () {
                if (this.keyName == "") {
                    return
                }
                this.cid="";
                pageNo++;
                this.$http.get(url+'search?key=' + this.keyName + '&pageNo=' + pageNo).then(function (res) {
                        if (pageNo == 1) {
                            this.sites = res.body.result_list.map_data;
                        } else {
                            this.sites = this.sites.concat(res.body.result_list.map_data);
                        }
                    }
                );
            },
            buyingSearch: function () {
                pageNo++;
                this.$http.get(url+'panic_buying2?key=生活用品&pageNo=' + pageNo).then(function (res) {
                    console.log(res.body.result_list.map_data)
                    if (pageNo == 1) {
                        this.sites = res.body.result_list.map_data;
                    } else {
                        this.sites = this.sites.concat(res.body.result_list.map_data);
                    }

                }, function (rsp) {
                    // console.log(rsp.status);
                });
            },
            copyCommand: function (message, addr) {
                this.hint = "复制口令";
                this.$http.get(url+'getTaoCommand?text=' + message + '&url=' + addr).then(function (res) {
                    this.cmd = res.body;
                        showBg();//显示复制口令窗口
                    }
                );
            },

            couponList: function () {
                if (this.commShow == false) {
                    this.buyShow=false;
                    this.commShow=true;
                    this.cid=0;
                    pageNo=0;
                    topUp();
                    this.defaultPage();
                }
            },
            firstPage:function(){
                this.buyShow=false;
                this.commShow=true;
                this.cid="";

                pageNo=0;
                topUp();
                this.defaultPage();
            },
            buying:function () {
                if (this.buyShow==false){
                    this.buyShow = true;
                    this.commShow=false;
                    this.cid=0;
                    pageNo=0;
                    topUp();
                    this.buyingSearch();
                }
            },
            cating:function(cat){
                this.buyShow = false;
                this.commShow=true;
                this.cid = cat;
                pageNo=0;
                topUp();
                this.getCatCommodidy()
            },
            topUp:function () {
                topUp()
            },
            //控制显示回到顶部图标
            showTop:function (res) {
                this.topShow=res;
            },

        }
    });

    //显示复制口令窗口
    function showBg() {
        $("#dialog").show();
    }
    //关闭复制口令窗口
    function closeBg() {
        $("#dialog").hide();
    }
    //回到顶部
    function topUp() {
        (function smoothscroll(){
            var currentScroll = document.documentElement.scrollTop || document.body.scrollTop;
            if (currentScroll > 0) {
                window.requestAnimationFrame(smoothscroll);
                window.scrollTo (0,currentScroll - (currentScroll/5));
            }
        })();
    }

   window.onload = function(){
       app.defaultPage();
   };

   window.onscroll = function(){
       var scroll = true;
       //变量scrollTop是滚动条滚动时，距离顶部的距离
       let scrollTop = document.documentElement.scrollTop||document.body.scrollTop  || window.pageYOffset;
       //变量windowHeight是可视区的高度
       let windowHeight = document.documentElement.clientHeight || document.body.clientHeight;
       //变量scrollHeight是滚动条的总高度
       let scrollHeight = document.documentElement.scrollHeight||document.body.scrollHeight;

       if(scrollTop >=  document.documentElement.clientHeight){
           app.showTop(true)
       }else {
           app.showTop(false)
       }

       // console.log("距离顶部的距离:"+scrollTop+",可视区的高度:"+windowHeight+"{"+(scrollTop+windowHeight)+"},滚动条的总高度:"+scrollHeight)
       if(scrollTop+windowHeight+50>=scrollHeight ) {
           // 写后台加载数据的函数
            if (loading == true){
                loading=false;
                // console.log("距离顶部的距离:"+scrollTop+",可视区的高度:"+windowHeight+"{"+(scrollTop+windowHeight)+"},滚动条的总高度:"+scrollHeight);
                if(app.commShow){
                    if (app.cid == ""){
                        if (app.keyName == "") {
                            app.defaultPage();
                        } else {
                            app.search();
                        }
                    }else{
                        app.getCatCommodidy()
                    }
                }else{
                    app.buyingSearch();
                }
            }
       }
   };

   var clipboard = new ClipboardJS('#cpy_cmd');
    clipboard.on('success', function(e) {
        app.hint="复制成功"
    });

    clipboard.on('error', function(e) {
        app.hint="复制失败"
    });
