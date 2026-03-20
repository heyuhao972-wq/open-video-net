function getQuery(name){

    const params = new URLSearchParams(window.location.search)

    return params.get(name)

}

const state = {
    mode: "list",
    page: 1,
    limit: 10,
    query: "",
    loading: false,
    hasMore: false
}

function getLoadMoreButton(){
    return document.getElementById("video-load-more")
}

function setLoadMoreVisible(visible){
    const btn = getLoadMoreButton()
    if (!btn){
        return
    }
    btn.style.display = visible ? "inline-block" : "none"
}

async function renderVideoURIs(uris, append){
    const list = document.getElementById("video-list")
    if (!list){
        return
    }
    if (!append){
        list.innerHTML = ""
    }
    for (const uri of uris){
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
        let stats = null
        let commentCount = null
        try {
            stats = await getVideoStats(info.platformId, info.videoId)
            commentCount = await getCommentCount(info.platformId, info.videoId)
        } catch (e) {
            // ignore
        }
        const li = document.createElement("li")
        const link = document.createElement("a")
        link.href = "player.html?uri=" + encodeURIComponent(uri)
        link.innerText = (detail && detail.title) ? detail.title : info.videoId
        li.appendChild(link)
        if (detail && detail.cover_path){
            const img = document.createElement("img")
            img.src = getPlatformBase(info.platformId) + "/" + detail.cover_path.replace(/^\.?[\\/]/,"")
            img.width = 120
            li.appendChild(img)
        }
        if (stats && !stats.error){
            const meta = document.createElement("div")
            const comments = commentCount && commentCount.count ? commentCount.count : 0
            meta.innerText = "views: " + (stats.watches || 0) + " | likes: " + (stats.likes || 0) + " | shares: " + (stats.shares || 0) + " | favorites: " + (stats.favorites || 0) + " | comments: " + comments
            li.appendChild(meta)
        }
        list.appendChild(li)
    }
}

async function renderVideoObjects(videos, append){
    const list = document.getElementById("video-list")
    if (!list){
        return
    }
    if (!append){
        list.innerHTML = ""
    }
    for (const v of videos){
        const li = document.createElement("li")
        const link = document.createElement("a")
        const platformId = v.platform_id || "platformA"
        link.href = "player.html?platform=" + encodeURIComponent(platformId) + "&id=" + v.id
        link.innerText = v.title
        li.appendChild(link)
        if (v.cover_path){
            const img = document.createElement("img")
            img.src = getPlatformBase(platformId) + "/" + v.cover_path.replace(/^\.?[\\/]/,"")
            img.width = 120
            li.appendChild(img)
        }
        try {
            const stats = await getVideoStats(platformId, v.id)
            const commentCount = await getCommentCount(platformId, v.id)
            if (stats && !stats.error){
                const meta = document.createElement("div")
                const comments = commentCount && commentCount.count ? commentCount.count : 0
                meta.innerText = "views: " + (stats.watches || 0) + " | likes: " + (stats.likes || 0) + " | shares: " + (stats.shares || 0) + " | favorites: " + (stats.favorites || 0) + " | comments: " + comments
                li.appendChild(meta)
            }
        } catch (e) {
            // ignore
        }
        list.appendChild(li)
    }
}

function initRecommendSelector(){
    const select = document.getElementById("recommend-select")
    const custom = document.getElementById("recommend-custom")
    const apply = document.getElementById("recommend-apply")
    if (!select || !apply){
        return
    }

    const options = getRecommendOptions()
    select.innerHTML = ""
    options.forEach(opt=>{
        const o = document.createElement("option")
        o.value = opt.url
        o.innerText = opt.label + " (" + opt.url + ")"
        select.appendChild(o)
    })

    const current = getRecommendBase()
    if (current){
        select.value = current
    }

    apply.addEventListener("click", ()=>{
        const customValue = custom ? custom.value.trim() : ""
        const value = customValue || select.value
        if (value){
            setRecommendBase(value)
            if (custom){
                custom.value = ""
            }
            resetToRecommend()
            loadVideos()
        }
    })
}

function resetToRecommend(){
    state.mode = "recommend"
    state.page = 1
    state.query = ""
    state.hasMore = false
}

function resetToSearch(query){
    state.mode = "search"
    state.page = 1
    state.query = query || ""
    state.hasMore = false
}

function resetToList(){
    state.mode = "list"
    state.page = 1
    state.query = ""
    state.hasMore = false
}

async function loadRecommendPage(){
    const user = getQuery("user") || "demo"
    const typeSelect = document.getElementById("recommend-type")
    const recType = typeSelect ? typeSelect.value : ""

    try {
        const recommend = await getRecommend(user, recType, state.page, state.limit)
        const videos = recommend && Array.isArray(recommend.videos) ? recommend.videos : []
        if (videos.length === 0){
            setLoadMoreVisible(false)
            return false
        }
        const list = document.getElementById("video-list")
        const beforeCount = list ? list.children.length : 0
        state.hasMore = videos.length === state.limit
        if (typeof videos[0] === "string"){
            await renderVideoURIs(videos, state.page > 1)
        } else {
            await renderVideoObjects(videos, state.page > 1)
        }
        const afterCount = list ? list.children.length : 0
        if (afterCount <= beforeCount) {
            setLoadMoreVisible(false)
            return false
        }
        setLoadMoreVisible(state.hasMore)
        return true
    } catch (e) {
        return false
    }
}

async function loadListPage(){
    const videos = await getVideos(state.page, state.limit)
    state.hasMore = videos.length === state.limit
    await renderVideoObjects(videos, state.page > 1)
    setLoadMoreVisible(state.hasMore)
}

async function loadSearchPage(){
    const videos = await searchVideos(state.query, state.page, state.limit)
    state.hasMore = videos.length === state.limit
    await renderVideoObjects(videos, state.page > 1)
    setLoadMoreVisible(state.hasMore)
}

async function loadVideos(){
    if (state.loading){
        return
    }
    state.loading = true
    try {
        if (state.mode === "recommend"){
            const ok = await loadRecommendPage()
            if (!ok){
                resetToList()
                await loadListPage()
            }
        } else if (state.mode === "search"){
            await loadSearchPage()
        } else {
            await loadListPage()
        }
    } finally {
        state.loading = false
    }
}

const searchForm = document.getElementById("search-form")
if (searchForm){
    searchForm.addEventListener("submit", async e=>{
        e.preventDefault()
        const q = document.getElementById("search-input").value
        resetToSearch(q)
        await loadVideos()
    })
}

const userSearchForm = document.getElementById("user-search-form")
if (userSearchForm){
    userSearchForm.addEventListener("submit", async e=>{
        e.preventDefault()
        const q = document.getElementById("user-search-input").value
        const res = await searchUsers(q)
        const list = document.getElementById("user-list")
        list.innerHTML = ""
        const users = res.users || []
        users.forEach(u=>{
            const li = document.createElement("li")
            const link = document.createElement("a")
            link.href = "user.html?id=" + encodeURIComponent(u.id)
            link.innerText = u.nickname || u.id
            li.appendChild(link)
            list.appendChild(li)
        })
    })
}

const tagSearchForm = document.getElementById("tag-search-form")
if (tagSearchForm){
    tagSearchForm.addEventListener("submit", async e=>{
        e.preventDefault()
        const tag = document.getElementById("tag-search-input").value.trim()
        if (!tag){
            return
        }
        window.location.href = "tags.html?tag=" + encodeURIComponent(tag)
    })
}

loadVideos()
initRecommendSelector()

const typeSelect = document.getElementById("recommend-type")
if (typeSelect){
    typeSelect.addEventListener("change", ()=>{
        resetToRecommend()
        loadVideos()
    })
}

const loadMoreBtn = getLoadMoreButton()
if (loadMoreBtn){
    loadMoreBtn.addEventListener("click", ()=>{
        if (!state.hasMore || state.loading){
            return
        }
        state.page += 1
        loadVideos()
    })
}
