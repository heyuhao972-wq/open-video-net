async function loadNotifications(){
    const list = document.getElementById("notification-list")
    if (!list){
        return
    }
    const res = await getNotifications()
    if (res.error){
        list.innerText = res.error
        return
    }
    list.innerHTML = ""
    const items = res.notifications || []
    items.forEach(n=>{
        const li = document.createElement("li")
        const text = document.createElement("span")
        const readFlag = n.read ? "[read] " : "[new] "
        text.innerText = readFlag + n.type + " from " + (n.actor_id || "") + " video " + (n.video_id || "")
        li.appendChild(text)

        const btn = document.createElement("button")
        btn.innerText = "Mark Read"
        btn.addEventListener("click", async ()=>{
            const r = await markNotificationsRead(n.id, false)
            if (r.error){
                alert(r.error)
                return
            }
            await loadNotifications()
        })
        li.appendChild(btn)

        list.appendChild(li)

        enrichNotification(n, text)
    })
}

async function enrichNotification(n, textEl){
    const parts = []
    if (n.actor_id){
        const actor = await getProfileById(n.actor_id)
        if (actor && actor.user && actor.user.nickname){
            parts.push("from " + actor.user.nickname)
        } else {
            parts.push("from " + n.actor_id)
        }
    }
    if (n.video_id){
        const video = await getVideoById(n.video_id)
        if (video && video.title){
            parts.push("video " + video.title)
        } else {
            parts.push("video " + n.video_id)
        }
    }
    const readFlag = n.read ? "[read] " : "[new] "
    const head = n.type || "event"
    if (parts.length > 0){
        textEl.innerText = readFlag + head + " " + parts.join(" ")
    }
}

const markAll = document.getElementById("mark-all")
if (markAll){
    markAll.addEventListener("click", async ()=>{
        const r = await markNotificationsRead(0, true)
        if (r.error){
            alert(r.error)
            return
        }
        await loadNotifications()
    })
}

loadNotifications()
