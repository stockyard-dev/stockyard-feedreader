package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Feedreader</title>
<link href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:ital,wght@0,400;0,700;1,400&family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:960px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center;font-family:var(--mono)}
.st-v{font-size:1.3rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center;flex-wrap:wrap}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.feed{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}
.feed:hover{border-color:var(--leather)}
.feed-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.feed-title{font-size:.9rem}.feed-title a{color:var(--cream);text-decoration:none}.feed-title a:hover{color:var(--rust)}
.feed-url{font-family:var(--mono);font-size:.6rem;color:var(--cm);margin-top:.15rem;word-break:break-all}
.feed-meta{font-family:var(--mono);font-size:.55rem;color:var(--cm);margin-top:.35rem;display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
.feed-actions{display:flex;gap:.3rem;flex-shrink:0}
.badge{font-family:var(--mono);font-size:.5rem;padding:.15rem .4rem;text-transform:uppercase;letter-spacing:1px;border:1px solid}
.badge.active{border-color:var(--green);color:var(--green)}.badge.paused{border-color:var(--gold);color:var(--gold)}.badge.error{border-color:var(--red);color:var(--red)}
.cat-badge{font-family:var(--mono);font-size:.5rem;padding:.1rem .35rem;background:var(--bg3);color:var(--cd)}
.unread-badge{font-family:var(--mono);font-size:.55rem;padding:.1rem .35rem;background:var(--rust);color:#fff;min-width:18px;text-align:center}
.btn{font-family:var(--mono);font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}.btn-p:hover{background:#d4682f}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}.fr label{display:block;font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(3,1fr)}.row2{grid-template-columns:1fr}.toolbar{flex-direction:column}.search{min-width:100%}}
</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> FEEDREADER</h1><button class="btn btn-p" onclick="openForm()">+ Add Feed</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search feeds..." oninput="render()">
<select class="filter-sel" id="cat-filter" onchange="render()"><option value="">All Categories</option></select>
<select class="filter-sel" id="status-filter" onchange="render()"><option value="">All Status</option><option value="active">Active</option><option value="paused">Paused</option><option value="error">Error</option></select>
</div>
<div id="feeds"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',feeds=[],editId=null;

async function load(){var r=await fetch(A+'/feeds').then(function(r){return r.json()});feeds=r.feeds||[];renderStats();buildCatFilter();render();}

function renderStats(){
var total=feeds.length;
var active=feeds.filter(function(f){return f.status==='active'}).length;
var totalUnread=feeds.reduce(function(s,f){return s+f.unread_count},0);
document.getElementById('stats').innerHTML=[
{l:'Feeds',v:total},{l:'Active',v:active},{l:'Unread',v:totalUnread,c:totalUnread>0?'var(--rust)':''}
].map(function(x){return '<div class="st"><div class="st-v" style="'+(x.c?'color:'+x.c:'')+'">'+x.v+'</div><div class="st-l">'+x.l+'</div></div>'}).join('');
}

function buildCatFilter(){
var cats={};feeds.forEach(function(f){if(f.category)cats[f.category]=true});
var sel=document.getElementById('cat-filter');var cur=sel.value;
sel.innerHTML='<option value="">All Categories</option>';
Object.keys(cats).sort().forEach(function(c){sel.innerHTML+='<option value="'+esc(c)+'"'+(cur===c?' selected':'')+'>'+esc(c)+'</option>';});
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var cat=document.getElementById('cat-filter').value;
var status=document.getElementById('status-filter').value;
var f=feeds;
if(cat)f=f.filter(function(x){return x.category===cat});
if(status)f=f.filter(function(x){return x.status===status});
if(q)f=f.filter(function(x){return(x.title||'').toLowerCase().includes(q)||(x.url||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('feeds').innerHTML='<div class="empty">No feeds found. Add your first RSS feed.</div>';return;}
var h='';f.forEach(function(i){
h+='<div class="feed"><div class="feed-top"><div style="flex:1">';
h+='<div class="feed-title">'+(i.site_url?'<a href="'+esc(i.site_url)+'" target="_blank" rel="noopener">'+esc(i.title)+' &#8599;</a>':esc(i.title))+'</div>';
h+='<div class="feed-url">'+esc(i.url)+'</div>';
h+='</div><div class="feed-actions">';
if(i.unread_count>0)h+='<span class="unread-badge">'+i.unread_count+'</span>';
h+='<button class="btn btn-sm" onclick="openEdit(\''+i.id+'\')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(\''+i.id+'\')" style="color:var(--red)">&#10005;</button>';
h+='</div></div>';
h+='<div class="feed-meta">';
h+='<span class="badge '+i.status+'">'+i.status+'</span>';
if(i.category)h+='<span class="cat-badge">'+esc(i.category)+'</span>';
h+='<span>'+i.item_count+' items</span>';
if(i.last_fetched_at)h+='<span>Fetched: '+ft(i.last_fetched_at)+'</span>';
h+='</div></div>';
});
document.getElementById('feeds').innerHTML=h;
}

async function del(id){if(!confirm('Remove this feed?'))return;await fetch(A+'/feeds/'+id,{method:'DELETE'});load();}

function formHTML(feed){
var i=feed||{title:'',url:'',site_url:'',category:'',status:'active'};
var isEdit=!!feed;
var h='<h2>'+(isEdit?'EDIT FEED':'ADD FEED')+'</h2>';
h+='<div class="fr"><label>Feed URL *</label><input id="f-url" value="'+esc(i.url)+'" placeholder="https://example.com/rss.xml"></div>';
h+='<div class="fr"><label>Title</label><input id="f-title" value="'+esc(i.title)+'" placeholder="Feed name (auto-detected on fetch)"></div>';
h+='<div class="row2"><div class="fr"><label>Site URL</label><input id="f-site" value="'+esc(i.site_url)+'" placeholder="https://example.com"></div>';
h+='<div class="fr"><label>Category</label><input id="f-cat" value="'+esc(i.category)+'" placeholder="e.g. tech, news"></div></div>';
if(isEdit){h+='<div class="fr"><label>Status</label><select id="f-status">';
['active','paused','error'].forEach(function(s){h+='<option value="'+s+'"'+(i.status===s?' selected':'')+'>'+s.charAt(0).toUpperCase()+s.slice(1)+'</option>';});
h+='</select></div>';}
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add Feed')+'</button></div>';
return h;
}

function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');document.getElementById('f-url').focus();}
function openEdit(id){var feed=null;for(var j=0;j<feeds.length;j++){if(feeds[j].id===id){feed=feeds[j];break;}}if(!feed)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(feed);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}

async function submit(){
var url=document.getElementById('f-url').value.trim();
if(!url){alert('Feed URL is required');return;}
var body={url:url,title:document.getElementById('f-title').value.trim()||url,site_url:document.getElementById('f-site').value.trim(),category:document.getElementById('f-cat').value.trim()};
if(editId){var sel=document.getElementById('f-status');if(sel)body.status=sel.value;
await fetch(A+'/feeds/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/feeds',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
closeModal();load();
}

function ft(t){if(!t)return'';try{return new Date(t).toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});
load();
</script></body></html>`
