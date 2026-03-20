const form = document.getElementById("tag-form")
const list = document.getElementById("tag-videos")
const loadMoreBtn = document.getElementById("tag-load-more")

const state = {
    tag: "",
    page: 1,
    limit: 10,
    hasMore: false,
    loading: false
}

function getQuery(name){
    const params = new URLSearchParams(window.location.search)
    return params.get(name)
}

function setLoadMoreVisible(visible){
    if (!loadMoreBtn){
        return
    }
    loadMoreBtn.style.display = visible ? "inline-block" : "none"
}

async function loadTag(tag, append){
    if (!list){
        return
    }
    const res = await searchVideosByTag(tag, state.page, state.limit)
    if (!append){
        list.innerHTML = ""
    }
    const videos = res.videos || []
    videos.forEach(v=>{
        const li = document.createElement("li")
        const link = document.createElement("a")
        const platformId = v.platform_id || "platformA"
        link.href = "player.html?platform=" + encodeURIComponent(platformId) + "&id=" + v.id
        link.innerText = v.title || v.id
        li.appendChild(link)
        if (v.cover_path){
            const img = document.createElement("img")
            img.src = getPlatformBase(platformId) + "/" + v.cover_path.replace(/^\.?[\\/]/,"")
            img.width = 120
            li.appendChild(img)
        }
        list.appendChild(li)
    })
    state.hasMore = videos.length === state.limit
    setLoadMoreVisible(state.hasMore)
}

if (form){
    form.addEventListener("submit", async e=>{
        e.preventDefault()
        const tag = document.getElementById("tag-input").value.trim()
        if (!tag){
            return
        }
        state.tag = tag
        state.page = 1
        await loadTag(tag, false)
    })
}

const initial = getQuery("tag")
if (initial){
    document.getElementById("tag-input").value = initial
    state.tag = initial
    state.page = 1
    loadTag(initial, false)
}

if (loadMoreBtn){
    loadMoreBtn.addEventListener("click", async ()=>{
        if (state.loading || !state.hasMore){
            return
        }
        state.loading = true
        try {
            state.page += 1
            await loadTag(state.tag, true)
        } finally {
            state.loading = false
        }
    })
}
