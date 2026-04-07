package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Feedreader</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5;font-size:13px}
.hdr{padding:.8rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.main{padding:1.2rem 1.5rem;max-width:1100px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700;color:var(--gold)}
.st-v.green{color:var(--green)}
.st-v.orange{color:var(--orange)}
.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}

.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(320px,1fr));gap:.6rem}
.card{background:var(--bg2);border:1px solid var(--bg3);padding:.9rem 1rem;display:flex;flex-direction:column;gap:.4rem;transition:border-color .15s}
.card:hover{border-color:var(--leather)}
.card.has-unread{border-left:3px solid var(--orange)}
.card.error{border-left:3px solid var(--red)}
.card.paused{opacity:.6}
.card.archived{opacity:.5}
.card-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem;cursor:pointer}
.card-title{font-size:.85rem;font-weight:700;color:var(--cream);flex:1}
.unread-badge{font-family:var(--mono);font-size:.7rem;font-weight:700;color:#fff;background:var(--orange);padding:.15rem .5rem;border-radius:10px;flex-shrink:0;line-height:1.2}
.unread-badge.zero{background:var(--bg3);color:var(--cm)}
.card-meta{font-size:.55rem;color:var(--cm);display:flex;gap:.5rem;flex-wrap:wrap;align-items:center}
.card-meta a{color:var(--cd);text-decoration:none}
.card-meta a:hover{color:var(--rust)}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm);font-weight:700}
.badge.active{border-color:var(--green);color:var(--green)}
.badge.paused{border-color:var(--cm);color:var(--cm)}
.badge.error{border-color:var(--red);color:var(--red)}
.badge.archived{border-color:var(--cm);color:var(--cm)}
.badge.cat{border-color:var(--leather);color:var(--leather)}
.card-counts{font-size:.6rem;color:var(--cd);font-family:var(--mono)}
.card-counts strong{color:var(--cream)}
.card-actions{display:flex;gap:.3rem;margin-top:.3rem}
.btn-read{font-family:var(--mono);font-size:.55rem;padding:.2rem .4rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--green);transition:.15s}
.btn-read:hover{border-color:var(--green)}
.btn-read:disabled{opacity:.4;cursor:not-allowed}
.card-extra{font-size:.55rem;color:var(--cd);margin-top:.4rem;padding-top:.3rem;border-top:1px dashed var(--bg3);display:flex;flex-direction:column;gap:.15rem}
.card-extra-row{display:flex;gap:.4rem}
.card-extra-label{color:var(--cm);text-transform:uppercase;letter-spacing:.5px;min-width:90px}
.card-extra-val{color:var(--cream)}

.btn{font-family:var(--mono);font-size:.6rem;padding:.3rem .55rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:.15s}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-p:hover{opacity:.85;color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.btn-del{color:var(--red);border-color:#3a1a1a}
.btn-del:hover{border-color:var(--red);color:var(--red)}

.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:520px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.fr-section{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.fr-section-label{font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.acts .btn-del{margin-right:auto}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(2,1fr)}}
</style>
</head>
<body>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> FEEDREADER</h1>
<button class="btn btn-p" onclick="openNew()">+ Add Feed</button>
</div>

<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search feeds..." oninput="debouncedRender()">
<select class="filter-sel" id="status-filter" onchange="render()">
<option value="">All Statuses</option>
<option value="active">Active</option>
<option value="paused">Paused</option>
<option value="error">Error</option>
<option value="archived">Archived</option>
</select>
<select class="filter-sel" id="category-filter" onchange="render()">
<option value="">All Categories</option>
</select>
<select class="filter-sel" id="unread-filter" onchange="render()">
<option value="">All Feeds</option>
<option value="unread">Unread Only</option>
</select>
</div>
<div id="grid" class="grid"></div>
</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var RESOURCE='feeds';

var fields=[
{name:'title',label:'Title',type:'text',required:true},
{name:'url',label:'Feed URL',type:'url',required:true,placeholder:'https://example.com/feed.xml'},
{name:'site_url',label:'Site URL',type:'url',placeholder:'https://example.com'},
{name:'category',label:'Category',type:'select_or_text',options:[]},
{name:'status',label:'Status',type:'select',options:['active','paused','error','archived']}
];

var feeds=[],feedExtras={},editId=null,searchTimer=null;

function fmtAgo(s){
if(!s)return'never';
try{
var d=new Date(s);
if(isNaN(d.getTime()))return s;
var diffMs=Date.now()-d;
if(diffMs<0)return'just now';
var sec=Math.floor(diffMs/1000);
if(sec<60)return sec+'s ago';
var min=Math.floor(sec/60);
if(min<60)return min+'m ago';
var hours=Math.floor(min/60);
if(hours<24)return hours+'h ago';
var days=Math.floor(hours/24);
if(days<7)return days+'d ago';
return d.toLocaleDateString('en-US',{month:'short',day:'numeric'});
}catch(e){return s}
}

function fieldByName(n){for(var i=0;i<fields.length;i++)if(fields[i].name===n)return fields[i];return null}

function debouncedRender(){
clearTimeout(searchTimer);
searchTimer=setTimeout(render,200);
}

async function load(){
try{
var resps=await Promise.all([
fetch(A+'/feeds').then(function(r){return r.json()}),
fetch(A+'/stats').then(function(r){return r.json()})
]);
feeds=resps[0].feeds||[];
renderStats(resps[1]||{});

try{
var ex=await fetch(A+'/extras/'+RESOURCE).then(function(r){return r.json()});
feedExtras=ex||{};
feeds.forEach(function(f){
var x=feedExtras[f.id];
if(!x)return;
Object.keys(x).forEach(function(k){if(f[k]===undefined)f[k]=x[k]});
});
}catch(e){feedExtras={}}

populateCategoryFilter();
}catch(e){
console.error('load failed',e);
feeds=[];
}
render();
}

function populateCategoryFilter(){
var sel=document.getElementById('category-filter');
if(!sel)return;
var current=sel.value;
var seen={};var cats=[];
feeds.forEach(function(f){if(f.category&&!seen[f.category]){seen[f.category]=true;cats.push(f.category)}});
cats.sort();
sel.innerHTML='<option value="">All Categories</option>'+cats.map(function(c){return'<option value="'+esc(c)+'"'+(c===current?' selected':'')+'>'+esc(c)+'</option>'}).join('');
}

function renderStats(s){
var total=s.total||0;
var totalUnread=s.total_unread||0;
var withUnread=s.with_unread||0;
var totalItems=s.total_items||0;
document.getElementById('stats').innerHTML=
'<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Feeds</div></div>'+
'<div class="st"><div class="st-v orange">'+totalUnread+'</div><div class="st-l">Unread Items</div></div>'+
'<div class="st"><div class="st-v">'+withUnread+'</div><div class="st-l">Feeds w/ Unread</div></div>'+
'<div class="st"><div class="st-v">'+totalItems+'</div><div class="st-l">Total Items</div></div>';
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var sf=document.getElementById('status-filter').value;
var cf=document.getElementById('category-filter').value;
var uf=document.getElementById('unread-filter').value;

var f=feeds.slice();
if(q)f=f.filter(function(fd){
return(fd.title||'').toLowerCase().includes(q)||
(fd.url||'').toLowerCase().includes(q)||
(fd.site_url||'').toLowerCase().includes(q);
});
if(sf)f=f.filter(function(fd){return fd.status===sf});
if(cf)f=f.filter(function(fd){return fd.category===cf});
if(uf==='unread')f=f.filter(function(fd){return(fd.unread_count||0)>0});

if(!f.length){
var msg=window._emptyMsg||'No feeds yet.';
document.getElementById('grid').innerHTML='<div class="empty" style="grid-column:1/-1">'+esc(msg)+'</div>';
return;
}

var h='';
f.forEach(function(fd){h+=cardHTML(fd)});
document.getElementById('grid').innerHTML=h;
}

function cardHTML(f){
var unread=f.unread_count||0;
var items=f.item_count||0;
var status=f.status||'active';
var cls='card '+status;
if(unread>0&&status!=='error')cls+=' has-unread';

var h='<div class="'+cls+'">';
h+='<div class="card-top" onclick="openEdit(\''+esc(f.id)+'\')">';
h+='<div class="card-title">'+esc(f.title)+'</div>';
h+='<span class="unread-badge'+(unread===0?' zero':'')+'">'+unread+'</span>';
h+='</div>';

h+='<div class="card-meta">';
h+='<span class="badge '+esc(status)+'">'+esc(status)+'</span>';
if(f.category)h+='<span class="badge cat">'+esc(f.category)+'</span>';
if(f.site_url)h+='<a href="'+esc(f.site_url)+'" target="_blank" rel="noopener" onclick="event.stopPropagation()">site</a>';
if(f.url)h+='<a href="'+esc(f.url)+'" target="_blank" rel="noopener" onclick="event.stopPropagation()">feed</a>';
h+='</div>';

h+='<div class="card-counts"><strong>'+items+'</strong> items &middot; fetched '+esc(fmtAgo(f.last_fetched_at))+'</div>';

h+='<div class="card-actions">';
h+='<button class="btn-read" onclick="markRead(\''+esc(f.id)+'\',event)"'+(unread===0?' disabled':'')+'>Mark all read</button>';
h+='</div>';

// Custom field display
var customRows='';
fields.forEach(function(fd){
if(!fd.isCustom)return;
var v=f[fd.name];
if(v===undefined||v===null||v==='')return;
customRows+='<div class="card-extra-row">';
customRows+='<span class="card-extra-label">'+esc(fd.label)+'</span>';
customRows+='<span class="card-extra-val">'+esc(String(v))+'</span>';
customRows+='</div>';
});
if(customRows)h+='<div class="card-extra">'+customRows+'</div>';

h+='</div>';
return h;
}

async function markRead(id,ev){
ev.stopPropagation();
try{
await fetch(A+'/feeds/'+id+'/read',{method:'POST'});
load();
}catch(e){alert('Failed')}
}

// ─── Modal ────────────────────────────────────────────────────────

function fieldHTML(f,value){
var v=value;
if(v===undefined||v===null)v='';
var req=f.required?' *':'';
var ph=f.placeholder?(' placeholder="'+esc(f.placeholder)+'"'):'';
var h='<div class="fr"><label>'+esc(f.label)+req+'</label>';

if(f.type==='select'){
h+='<select id="f-'+f.name+'">';
if(!f.required)h+='<option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=(String(v)===String(o))?' selected':'';
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(String(o))+'</option>';
});
h+='</select>';
}else if(f.type==='select_or_text'){
h+='<input list="dl-'+f.name+'" type="text" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
h+='<datalist id="dl-'+f.name+'">';
var opts=(f.options||[]).slice();
feeds.forEach(function(fd){if(fd.category&&opts.indexOf(fd.category)===-1)opts.push(fd.category)});
opts.forEach(function(o){h+='<option value="'+esc(String(o))+'">'});
h+='</datalist>';
}else if(f.type==='textarea'){
h+='<textarea id="f-'+f.name+'" rows="3"'+ph+'>'+esc(String(v))+'</textarea>';
}else if(f.type==='number'){
h+='<input type="number" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}else{
var inputType=f.type||'text';
h+='<input type="'+esc(inputType)+'" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}
h+='</div>';
return h;
}

function formHTML(feed){
var f=feed||{};
var isEdit=!!feed;
var h='<h2>'+(isEdit?'EDIT FEED':'NEW FEED')+'</h2>';

h+=fieldHTML(fieldByName('title'),f.title);
h+=fieldHTML(fieldByName('url'),f.url);
h+=fieldHTML(fieldByName('site_url'),f.site_url);
h+='<div class="row2">'+fieldHTML(fieldByName('category'),f.category)+fieldHTML(fieldByName('status'),f.status||'active')+'</div>';

var customFields=fields.filter(function(f){return f.isCustom});
if(customFields.length){
var label=window._customSectionLabel||'Additional Details';
h+='<div class="fr-section"><div class="fr-section-label">'+esc(label)+'</div>';
customFields.forEach(function(f){h+=fieldHTML(f,feed?feed[f.name]:'')});
h+='</div>';
}

h+='<div class="acts">';
if(isEdit)h+='<button class="btn btn-del" onclick="delItem()">Delete</button>';
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add')+'</button>';
h+='</div>';
return h;
}

function openNew(){
editId=null;
document.getElementById('mdl').innerHTML=formHTML();
document.getElementById('mbg').classList.add('open');
var n=document.getElementById('f-title');if(n)n.focus();
}

function openEdit(id){
var f=null;
for(var i=0;i<feeds.length;i++){if(feeds[i].id===id){f=feeds[i];break}}
if(!f)return;
editId=id;
document.getElementById('mdl').innerHTML=formHTML(f);
document.getElementById('mbg').classList.add('open');
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
editId=null;
}

async function submit(){
var titleEl=document.getElementById('f-title');
if(!titleEl||!titleEl.value.trim()){alert('Title is required');return}
var urlEl=document.getElementById('f-url');
if(!urlEl||!urlEl.value.trim()){alert('URL is required');return}

var body={};
var extras={};
fields.forEach(function(f){
var el=document.getElementById('f-'+f.name);
if(!el)return;
var val=el.value.trim();
if(f.isCustom)extras[f.name]=val;
else body[f.name]=val;
});

var savedId=editId;
try{
if(editId){
var r1=await fetch(A+'/feeds/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/feeds',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Add failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/'+RESOURCE+'/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){alert('Network error: '+e.message);return}
closeModal();
load();
}

async function delItem(){
if(!editId)return;
if(!confirm('Delete this feed?'))return;
await fetch(A+'/feeds/'+editId,{method:'DELETE'});
closeModal();
load();
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

// Auto-refresh every 60s
setInterval(load,60000);

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;
if(cfg.primary_label)window._customSectionLabel=cfg.primary_label+' Details';

if(Array.isArray(cfg.categories)){
var catField=fieldByName('category');
if(catField)catField.options=cfg.categories;
}

if(Array.isArray(cfg.custom_fields)){
cfg.custom_fields.forEach(function(cf){
if(!cf||!cf.name||!cf.label)return;
if(fieldByName(cf.name))return;
fields.push({
name:cf.name,
label:cf.label,
type:cf.type||'text',
options:cf.options||[],
isCustom:true
});
});
}
}).catch(function(){
}).finally(function(){
load();
});
})();
</script>
</body>
</html>`
