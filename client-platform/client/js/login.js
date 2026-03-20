const loginForm = document.getElementById("login-form")
const keyJson = document.getElementById("key-json")
const importBtn = document.getElementById("import-keys")
const exportBtn = document.getElementById("export-keys")
const downloadBtn = document.getElementById("download-keys")

loginForm.addEventListener("submit", async e=>{

    e.preventDefault()

    const inputKey = document.getElementById("public_key").value.trim()
    const publicKey = inputKey || localStorage.getItem("public_key")
    if (!publicKey){
        alert("public key required")
        return
    }
    const res = await loginUser(publicKey)

    if (res.error){
        alert(res.error)
        return
    }

    localStorage.setItem("token", res.token)
    localStorage.setItem("user_id", res.user.id)

    window.location.href = "index.html"

})

importBtn.addEventListener("click", ()=>{
    const res = importKeyBundle(keyJson.value)
    if (res.error){
        alert(res.error)
        return
    }
    alert("keys imported")
})

exportBtn.addEventListener("click", ()=>{
    const payload = exportKeyBundle()
    if (!payload){
        alert("no keys in local storage")
        return
    }
    keyJson.value = payload
})

downloadBtn.addEventListener("click", ()=>{
    const res = downloadKeyBundle()
    if (res.error){
        alert(res.error)
    }
})
