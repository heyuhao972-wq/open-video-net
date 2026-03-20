function getQuery(name){

    const params = new URLSearchParams(window.location.search)

    return params.get(name)

}

async function loadVideos(){

    const list = document.getElementById("video-list")
    const recommendSelect = document.getElementById("recommend-select")
    if (recommendSelect){
        const saved = localStorage.getItem("recommend_base")
        if (saved){
            recommendSelect.value = saved
        }
        recommendSelect.addEventListener("change", ()=>{
            setRecommendBase(recommendSelect.value)
            loadVideos()
        })
    }

    const user = getQuery("user") || "demo"

    let videos = []

    try {
        const recommend = await getRecommend(user)
        if (recommend && Array.isArray(recommend.videos) && recommend.videos.length > 0) {
            videos = recommend.videos
        }
    } catch (e) {
        // fallback to video list
    }

    if (videos.length === 0) {
        videos = await getVideos()
    }

    list.innerHTML = ""

    if (videos.length > 0 && typeof videos[0] === "string"){
        for (const uri of videos){
            const info = parseVideoURI(uri)
            if (!info){
                continue
            }
            let detail = null
            try {
                detail = await getVideoByPlatform(info.platformId, info.videoId)
            } catch (e) {
                // ignore
            }
            const card = buildCard({
                id: info.videoId,
                title: (detail && detail.title) ? detail.title : info.videoId,
                tags: detail && detail.tags ? detail.tags : [],
                platform_id: info.platformId
            })
            card.querySelector("a").href = "player.html?uri=" + encodeURIComponent(uri)
            list.appendChild(card)
        }
        return
    }

    videos.forEach(v=>{
        const card = buildCard(v)
        list.appendChild(card)
    })

}

    const searchForm = document.getElementById("search-form")
if (searchForm){
    searchForm.addEventListener("submit", async e=>{
        e.preventDefault()
        const q = document.getElementById("search-input").value
        const results = await searchVideos(q)
        const list = document.getElementById("video-list")
        list.innerHTML = ""
        results.forEach(v=>{
            const card = buildCard(v)
            list.appendChild(card)
        })
    })
}

loadVideos()

function buildCard(v){
    const card = document.createElement("div")
    card.className = "video-card"

    const thumb = document.createElement("div")
    thumb.className = "video-thumb"
    card.appendChild(thumb)

    const link = document.createElement("a")
    link.className = "video-title"
    link.innerText = v.title || v.id
    const platformId = v.platform_id || "platformA"
    link.href = "player.html?platform=" + encodeURIComponent(platformId) + "&id=" + v.id
    card.appendChild(link)

    const meta = document.createElement("div")
    meta.className = "video-meta"
    meta.innerText = platformId + " • " + (v.views || 0) + " views"
    card.appendChild(meta)

    if (Array.isArray(v.tags) && v.tags.length > 0){
        const tag = document.createElement("div")
        tag.className = "tag"
        tag.innerText = v.tags[0]
        card.appendChild(tag)
    }

    return card
}
