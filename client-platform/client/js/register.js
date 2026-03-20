const registerForm = document.getElementById("register-form")
const keyJson = document.getElementById("key-json")
const importBtn = document.getElementById("import-keys")
const exportBtn = document.getElementById("export-keys")
const downloadBtn = document.getElementById("download-keys")
const generateBtn = document.getElementById("generate-keys")

registerForm.addEventListener("submit", async e=>{

    e.preventDefault()

    const inputKey = document.getElementById("public_key").value.trim()
    let publicKey = inputKey
    let privateKey = localStorage.getItem("private_key")
    if (!publicKey){
        if (!privateKey || !localStorage.getItem("public_key")){
            const kp = await generateKeypair()
            publicKey = kp.public_key
            privateKey = kp.private_key
            localStorage.setItem("public_key", publicKey)
            localStorage.setItem("private_key", privateKey)
        } else {
            publicKey = localStorage.getItem("public_key")
        }
    }

    const res = await registerUser(publicKey)

    if (res.error){
        alert(res.error)
        return
    }

    if (res.user && res.user.public_key){
        localStorage.setItem("public_key", res.user.public_key)
    }

    alert("register success")
    window.location.href = "login.html"

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

generateBtn.addEventListener("click", async ()=>{
    const kp = await generateKeypair()
    localStorage.setItem("public_key", kp.public_key)
    localStorage.setItem("private_key", kp.private_key)
    keyJson.value = JSON.stringify(kp)
    const input = document.getElementById("public_key")
    if (input){
        input.value = kp.public_key
    }
    alert("keys generated")
})
