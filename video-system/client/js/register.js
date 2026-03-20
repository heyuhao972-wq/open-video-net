const registerForm = document.getElementById("register-form")
const exportBtn = document.getElementById("export-keys")
const importInput = document.getElementById("import-keys")
const publicKeyBox = document.getElementById("public-key")

registerForm.addEventListener("submit", async e=>{

    e.preventDefault()

    const username = document.getElementById("username").value

    const keys = await generateKeypair()
    localStorage.setItem("private_key", keys.privateKey)
    localStorage.setItem("public_key", keys.publicKey)
    if (username){
        localStorage.setItem("username", username)
    }
    if (publicKeyBox){
        publicKeyBox.value = keys.publicKey
    }

    const res = await registerUser(keys.publicKey)

    if (res.error){
        alert(res.error)
        return
    }

    if (res.user && res.user.id){
        localStorage.setItem("user_id", res.user.id)
    }

    alert("register success, private key saved in local storage")
    window.location.href = "login.html"

})

if (exportBtn){
    exportBtn.addEventListener("click", ()=>{
        const publicKey = localStorage.getItem("public_key") || ""
        const privateKey = localStorage.getItem("private_key") || ""
        if (!publicKey || !privateKey){
            alert("no keys found, please register first")
            return
        }
        const payload = JSON.stringify({public_key: publicKey, private_key: privateKey}, null, 2)
        const blob = new Blob([payload], {type: "application/json"})
        const url = URL.createObjectURL(blob)
        const a = document.createElement("a")
        a.href = url
        a.download = "openvideo-keys.json"
        a.click()
        URL.revokeObjectURL(url)
    })
}

if (importInput){
    importInput.addEventListener("change", async e=>{
        const file = e.target.files[0]
        if (!file){
            return
        }
        const text = await file.text()
        try {
            const data = JSON.parse(text)
            if (!data.public_key || !data.private_key){
                alert("invalid key file")
                return
            }
            localStorage.setItem("public_key", data.public_key)
            localStorage.setItem("private_key", data.private_key)
            if (publicKeyBox){
                publicKeyBox.value = data.public_key
            }
            alert("keys imported")
        } catch (err) {
            alert("invalid key file")
        }
    })
}
