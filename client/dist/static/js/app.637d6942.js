(function(e){function t(t){for(var o,a,c=t[0],r=t[1],l=t[2],u=0,p=[];u<c.length;u++)a=c[u],Object.prototype.hasOwnProperty.call(i,a)&&i[a]&&p.push(i[a][0]),i[a]=0;for(o in r)Object.prototype.hasOwnProperty.call(r,o)&&(e[o]=r[o]);d&&d(t);while(p.length)p.shift()();return s.push.apply(s,l||[]),n()}function n(){for(var e,t=0;t<s.length;t++){for(var n=s[t],o=!0,c=1;c<n.length;c++){var r=n[c];0!==i[r]&&(o=!1)}o&&(s.splice(t--,1),e=a(a.s=n[0]))}return e}var o={},i={app:0},s=[];function a(t){if(o[t])return o[t].exports;var n=o[t]={i:t,l:!1,exports:{}};return e[t].call(n.exports,n,n.exports,a),n.l=!0,n.exports}a.m=e,a.c=o,a.d=function(e,t,n){a.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:n})},a.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},a.t=function(e,t){if(1&t&&(e=a(e)),8&t)return e;if(4&t&&"object"===typeof e&&e&&e.__esModule)return e;var n=Object.create(null);if(a.r(n),Object.defineProperty(n,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var o in e)a.d(n,o,function(t){return e[t]}.bind(null,o));return n},a.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return a.d(t,"a",t),t},a.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},a.p="";var c=window["webpackJsonp"]=window["webpackJsonp"]||[],r=c.push.bind(c);c.push=t,c=c.slice();for(var l=0;l<c.length;l++)t(c[l]);var d=r;s.push([0,"chunk-vendors"]),n()})({0:function(e,t,n){e.exports=n("56d7")},"034f":function(e,t,n){"use strict";var o=n("85ec"),i=n.n(o);i.a},"31f2":function(e,t,n){"use strict";var o=n("ebce"),i=n.n(o);i.a},"56d7":function(e,t,n){"use strict";n.r(t);n("e260"),n("e6cf"),n("cca6"),n("a79d");var o=n("2b0e"),i=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{attrs:{id:"app"}},[n("HelloWorld",{attrs:{msg:"Welcome to Your Vue.js App"}})],1)},s=[],a=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("el-container",[n("div",{staticClass:"demo-input-size"},[e._v(" id:"),n("el-input",{attrs:{placeholder:"请输入设备ID"},model:{value:e.clientId,callback:function(t){e.clientId=t},expression:"clientId"}}),e._v(" topic:"),n("el-input",{attrs:{placeholder:"请输入主题"},model:{value:e.topic,callback:function(t){e.topic=t},expression:"topic"}}),n("el-button",{attrs:{disabled:e.StartDisabled},on:{click:e.start}},[e._v("连接")]),n("el-tag",{attrs:{type:"success"}},[e._v(e._s(e.iceConnectionState))])],1),n("div",[n("div",[n("el-input",{attrs:{placeholder:"请输入地址,默认为截屏"},model:{value:e.shotUrl,callback:function(t){e.shotUrl=t},expression:"shotUrl"}}),n("el-button",{on:{click:function(t){return e.job("shot")}}},[e._v("点击截图")]),n("el-dropdown",[n("el-button",[e._v(" "+e._s(e.areaTag)),n("i",{staticClass:"el-icon-arrow-down el-icon--right"})]),n("el-dropdown-menu",{attrs:{slot:"dropdown"},slot:"dropdown"},e._l(e.areainfos,(function(t,o){return n("el-dropdown-item",{key:o},[n("span",{on:{click:function(t){return e.shotByArea(o)}}},[e._v(e._s(t))])])})),1)],1)],1),n("br"),n("div",{staticStyle:{width:"200px"}},[n("el-image",{staticStyle:{height:"200px"},attrs:{src:e.url,"preview-src-list":e.srcList}},[n("div",{staticClass:"image-slot",attrs:{slot:"error"},slot:"error"},[n("i",{staticClass:"el-icon-picture-outline"})])]),n("el-progress",{attrs:{"text-inside":!1,"stroke-width":5,percentage:e.up_percent,status:e.loadSuccessful}})],1)]),n("div",{attrs:{id:"playsound"}},[n("div",[n("el-input",{attrs:{placeholder:"请输入mp3地址"},model:{value:e.soundUrl,callback:function(t){e.soundUrl=t},expression:"soundUrl"}}),n("el-button",{on:{click:function(t){return e.job("playsound")}}},[e._v("远程播放")])],1),n("br")]),n("div",[n("el-input",{attrs:{placeholder:"请输入视频地址"},model:{value:e.videoUrl,callback:function(t){e.videoUrl=t},expression:"videoUrl"}}),n("el-button",{on:{click:function(t){return e.job("live")}}},[e._v("开启远程视频")]),n("div",{attrs:{id:"remoteVideos",width:"400",height:"400"}}),n("video",{attrs:{id:"video1",width:"200",height:"200"}}),n("video",{attrs:{id:"video2",width:"200",height:"200"}})],1)])},c=[],r=(n("b0c0"),n("b680"),n("d3b7"),n("ac1f"),n("3ca3"),n("1276"),n("ddb0"),n("2b3d"),{data:function(){return{apiUrl:"http://localhost:9090/message",pc:null,ctlChannel:null,sdphannel:null,pcConfig:{iceServers:[{urls:"stun:47.99.78.179:3478"}],sdpSemantics:"unified-plan"},exchange_offer:"",exchange_answer:"",latesOffer:"",answerSd:"",logstatus:"",clientId:"",topic:"",up_percent:0,shotUrl:"",soundUrl:"./Lame_Drivers_-_01_-_Frozen_Egg.mp3",videoUrl:"rtsp://admin:admin123@192.168.2.241/h264/ch1/main/av_stream",url:"",srcList:[],localStream:null,remoteStream:null,firstIcecandidate:!0,activeVideos:0,totalVodeos:0,areaTag:"区域地点",areas:["dsnpayKxRDB4tuWS","52Z4UmYjjkz30mdP","uUJ2-6ufUiTXSsG8","NZm6Bq3HaZAbnh5L","PbNO6q8S8eg6oZo2","q_UMJuuV3mf7nwVP","wxRaecxm_KeHGbvq","mXyb6f6pyYbcnp2x","i_SNJngBNfLiMZ9y","ZD5S7yhnmsiSJ0Fn","pSoYnF3iNapnJ9II","x5Kkk1PuHKOIQAsI","y11q1GdIpzLDE1t_","RdBL8_FfdOJH2UUX"],areainfos:["遂宁市射洪县青岗镇","绵阳市游仙区朝真乡","绵阳市游仙区忠兴镇","绵阳市江油市小溪坝镇","绵阳市江油市大堰镇","德阳市罗江区略坪镇","德阳市罗江区新盛镇","泸州市泸县百和镇","泸州市泸县奇峰镇","成都市大邑县王泗镇","眉山市彭山县义和乡","眉山市彭山县公义镇","绵阳市梓潼县卧龙镇","绵阳市梓潼县黎雅镇"]}},watch:{answerSd:function(){this.startSession()}},computed:{StartDisabled:function(){return""==this.clientId||""==this.topic||"connected"==this.logstatus},ShottDisabled:function(){return"connected"!=this.logstatus},loadSuccessful:function(){return 100==this.up_percent?"success":"exception"},iceConnectionState:function(){return"lost"==this.logstatus?"连接失败":"disconnected"==this.logstatus?"连接断开":"checking"==this.logstatus?"连接中...":"connected"==this.logstatus?"连接成功":"准备连接"}},methods:{createPeerConnection:function(){var e=this;if(this.pc)console.warning("the pc have be created!");else{console.log("creat connecte"),this.pc=new RTCPeerConnection(this.pcConfig);var t=this.pc;t.oniceconnectionstatechange=function(){e.logstatus=t.iceConnectionState};var n=this;t.ontrack=function(e){var t=n.pc.getTransceivers();document.getElementById("remoteVideos").innerHTML="";for(var o=0;o<t.length;o++){var i=document.createElement(e.track.kind),s=t[o].receiver.track,a=new MediaStream;a.addTrack(s),i.setAttribute("width",200),i.setAttribute("height",200),i.srcObject=null,i.srcObject=a,i.autoplay=!0,i.controls=!0,document.getElementById("remoteVideos").appendChild(i)}},t.onnegotiationneeded=function(){console.log("onnegotiationneeded")},t.onicecandidate=function(){if(console.log("onicecandidate"),e.firstIcecandidate){console.log("onicecandidate***"),e.firstIcecandidate=!1;var o=btoa(JSON.stringify(t.localDescription));console.log("offerBase:",o);var i=e.clientId,s=e.topic,a="id="+i+"&topic="+s+"&type=4&offer="+o,c={method:"POST",headers:{Accept:"application/json","Content-Type":"application/x-www-form-urlencoded;charset=UTF-8"},body:a};console.log("onicecandidate1"),fetch(e.apiUrl,c).then((function(e){console.log("onicecandidate2"),e.ok?(console.log(e),e.json().then((function(e){var t=e.message.split(":")[1];console.log("answerSd",t),n.answerSd=t}))):console.log(e.status)})).catch((function(e){console.log(e.status)}))}},this.ctlChannel=t.createDataChannel("control");var o=this.ctlChannel;t.createOffer().then((function(e){return t.setLocalDescription(e)})),o.onopen=function(){console.log("ctlChannel is on")},o.onmessage=function(e){console.log(e)},t.ondatachannel=function(n){var o,i=[],s=0,a=n.channel;e.up_percent=0,a.onmessage=function(n){if("sdp"==a.label){if("ready"==n.data)return console.log("b0",t.getTransceivers().length),e.pc.addTransceiver("video"),console.log("b1",t.getTransceivers().length),console.log("b2",t.getTransceivers().length),void e.pc.createOffer().then((function(n){console.log("created:",n),console.log(e.pc),e.pc.setLocalDescription(n),console.log(n),a.send(btoa(JSON.stringify(n))),console.log("b3",t.getTransceivers().length)}));console.log("b4",t.getTransceivers().length);var c=new RTCSessionDescription(JSON.parse(atob(n.data)));console.log(c),e.pc.setRemoteDescription(c),console.log("b5",t.getTransceivers().length)}if("shot"==a.label)if(console.log("shot"),"string"==typeof n.data){var r=JSON.parse(atob(n.data));console.log(atob(n.data)),o=r.filesize}else if(i.push(n.data),s+=n.data.byteLength,e.up_percent=1*(s/o*100).toFixed(0),s===o){var l=new Blob(i);i=[],e.url=URL.createObjectURL(l),e.srcList.push(e.url)}}}}},start:function(){this.pc=null,this.createPeerConnection()},startSession:function(){var e=this.answerSd;if(""===e)return this.logstatus="lost",alert("连接失败");try{this.pc.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(e))))}catch(t){alert(t)}},job:function(e){var t=null,n=(new Date).getTime();switch(e){case"shot":this.url="";var o=this.shotUrl;""==o&&(o="screen"),t={id:this.clientId,topic:this.topic,type:"3",message:{shot:o,time:n}},t=btoa(JSON.stringify(t)),this.ctlChannel.send(t),this.sendCmd({});break;case"playsound":var i=this.soundUrl;null!=this.ctlChannel&&"open"==this.ctlChannel.readyState?(""==i?alert("mp3地址不能为空"):t={id:this.clientId,topic:this.topic,type:"1",message:{loop:1,path:this.soundUrl,time:n}},t=btoa(JSON.stringify(t)),this.ctlChannel.send(t)):this.sendCmd({category:1,times:1,path:i},(function(e){e.json().then((function(e){e=e.message.split(":")[1],console.log(e)}))}));break;case"live":if(""==this.videoUrl)return void alert("视频地址不能为空");t={url:this.videoUrl,type:"5"},t=btoa(JSON.stringify(t)),this.ctlChannel.send(t);break}},shotByArea:function(e){console.log("shotByArea"),this.shotUrl="",this.areaTag=this.areainfos[e];var t="code="+this.areas[e],n=this,o={method:"POST",headers:{Accept:"application/json","Content-Type":"application/x-www-form-urlencoded;charset=UTF-8"},body:t};fetch("http://118.122.120.57:18080/api/gis/video/live-url",o).then((function(e){e.ok?e.json().then((function(e){90010==e.code&&(console.log(e.data.list),n.shotUrl=e.data.list,n.job("shot"))})):console.log(e.status)})).catch((function(e){console.log(e.status)}))},startLocalStream:function(){navigator.getUserMedia=navigator.getUserMedia||navigator.webkitGetUserMedia||navigator.mozGetUserMedia;var e=this;navigator.getUserMedia?navigator.getUserMedia({audio:!0,video:{width:400,height:400}},(function(t){var n=document.querySelector("video");n.srcObject=t,e.localStream=t,console.log(t),n.onloadedmetadata=function(){n.play()}}),(function(e){console.log("The following error occurred: "+e.name)})):console.log("getUserMedia not supported")},sendCmd:function(e,t){if("connected"!=this.logstatus){var n=this.clientId,o=this.topic;""!=n&&""!=o||alert("id和topic不能为空");var i="id="+n+"&topic="+o+"&type="+e.category+"&times="+e.times+"&path="+e.path,s={method:"POST",headers:{Accept:"application/json","Content-Type":"application/x-www-form-urlencoded;charset=UTF-8"},body:i};fetch(this.apiUrl,s).then((function(e){e.ok?t(e):console.log(e.status)})).catch((function(e){console.log(e.status)}))}}}}),l=r,d=(n("31f2"),n("2877")),u=Object(d["a"])(l,a,c,!1,null,null,null),p=u.exports,h={name:"app",components:{HelloWorld:p}},f=h,g=(n("034f"),Object(d["a"])(f,i,s,!1,null,null,null)),v=g.exports,b=n("5c96"),m=n.n(b);n("0fae");o["default"].use(m.a),o["default"].config.productionTip=!1,new o["default"]({render:function(e){return e(v)}}).$mount("#app")},"85ec":function(e,t,n){},ebce:function(e,t,n){}});
//# sourceMappingURL=app.637d6942.js.map