package compilehtml

import (
	"encoding/json"
	"strconv"
	"strings"
)

func writeVariantSwitcher(w *strings.Builder, labels []string) {
	if len(labels) == 0 {
		return
	}
	labelsJSON, _ := json.Marshal(labels)

	w.WriteString(`<div id="psrt-variant-hint" class="psrt-variant-hint" aria-live="polite"></div>
<script>
(function(){
var labels=`)
	w.Write(labelsJSON)
	w.WriteString(`;
labels.push("");
var idx=0;
function apply(){
document.querySelectorAll(".psrt-text").forEach(function(el){
el.classList.add("psrt-hidden");
});
if(idx<labels.length-1){
document.querySelectorAll(".psrt-v-"+idx).forEach(function(el){
el.classList.remove("psrt-hidden");
});
}
var hint=document.getElementById("psrt-variant-hint");
if(hint){
hint.textContent=labels[idx]===""?"Sem PSRT (Ctrl+L)":labels[idx]+" (Ctrl+L)";
}
}
document.addEventListener("keydown",function(e){
if(e.ctrlKey&&!e.altKey&&e.key.toLowerCase()==="l"){
e.preventDefault();
idx=(idx+1)%labels.length;
apply();
}
});
apply();
})();
</script>
`)
}

func variantSwitcherCSS() string {
	return `
.psrt-hidden{display:none!important}
.psrt-variant-hint{position:fixed;z-index:9999;right:12px;bottom:12px;padding:6px 10px;font:12px/1.3 system-ui,sans-serif;color:#e8e8e8;background:rgba(0,0,0,.72);border-radius:6px;pointer-events:none;user-select:none}
`
}

func variantClass(v int) string {
	return "psrt-v-" + strconv.Itoa(v)
}
