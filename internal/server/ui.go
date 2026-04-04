package server
import "net/http"
func(s *Server)dashboard(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","text/html");w.Write([]byte(dashHTML))}
const dashHTML=`<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Feedreader</title>
<style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.hdr{padding:.8rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.layout{display:grid;grid-template-columns:220px 1fr;height:calc(100vh - 52px)}
@media(max-width:600px){.layout{grid-template-columns:1fr}}
.sidebar{border-right:1px solid var(--bg3);padding:.8rem;overflow-y:auto}
.feed{padding:.4rem .6rem;cursor:pointer;font-family:var(--mono);font-size:.72rem;color:var(--cd);margin-bottom:.1rem;display:flex;justify-content:space-between}
.feed:hover{background:var(--bg2);color:var(--cream)}.feed.active{color:var(--rust);background:var(--bg2)}
.feed-count{background:var(--bg3);padding:0 .3rem;font-size:.55rem;color:var(--cm);border-radius:2px}
.articles{padding:1rem;overflow-y:auto}
.article{border-bottom:1px solid var(--bg3);padding:.8rem 0;cursor:pointer}
.article:hover{background:var(--bg2);margin:0 -.5rem;padding:.8rem .5rem}
.article-title{font-size:.92rem;margin-bottom:.2rem}.article-title.unread{color:var(--cream);font-weight:700}.article-title.read{color:var(--cd)}
.article-meta{font-family:var(--mono);font-size:.6rem;color:var(--cm)}
.article-snippet{font-size:.78rem;color:var(--cm);margin-top:.2rem;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;overflow:hidden}
.btn{font-family:var(--mono);font-size:.6rem;padding:.25rem .6rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:400px;max-width:90vw}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.5rem}.fr label{display:block;font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.15rem}
.fr input{width:100%;padding:.35rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:.8rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.8rem}
</style></head><body>
<div class="hdr"><h1>FEEDREADER</h1><button class="btn btn-p" onclick="openForm()">+ Subscribe</button></div>
<div class="layout">
<div class="sidebar"><div style="font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem">Feeds</div><div id="feeds"></div></div>
<div class="articles" id="articles"><div class="empty">Subscribe to an RSS feed to get started</div></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api';let feeds=[],curFeed='';
async function load(){const r=await fetch(A+'/feeds').then(r=>r.json());feeds=r.feeds||[];renderFeeds();}
function renderFeeds(){let h='<div class="feed'+(curFeed===''?' active':'')+'" onclick="selectFeed(\'\')">All feeds <span class="feed-count">'+feeds.reduce((s,f)=>s+f.unread_count,0)+'</span></div>';
feeds.forEach(f=>{h+='<div class="feed'+(curFeed===f.id?' active':'')+'" onclick="selectFeed(\''+f.id+'\')">'+esc(f.title||f.url)+' <span class="feed-count">'+f.unread_count+'</span></div>';});
document.getElementById('feeds').innerHTML=h;}
function selectFeed(id){curFeed=id;renderFeeds();document.getElementById('articles').innerHTML='<div class="empty">Feed items will appear here after the first fetch.</div>';}
async function del(id){if(confirm('Unsubscribe?')){await fetch(A+'/feeds/'+id,{method:'DELETE'});load();}}
function openForm(){document.getElementById('mdl').innerHTML='<h2>Subscribe to Feed</h2><div class="fr"><label>RSS/Atom URL</label><input id="f-u" placeholder="https://example.com/feed.xml"></div><div class="fr"><label>Title (optional)</label><input id="f-t" placeholder="auto-detected from feed"></div><div class="fr"><label>Category</label><input id="f-c" placeholder="e.g. tech, news"></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="sub()">Subscribe</button></div>';document.getElementById('mbg').classList.add('open');}
async function sub(){await fetch(A+'/feeds',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({url:document.getElementById('f-u').value,title:document.getElementById('f-t').value,category:document.getElementById('f-c').value})});cm();load();}
function cm(){document.getElementById('mbg').classList.remove('open');}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();
</script></body></html>`
