//  Collections
//  ========================================== 
function listCollections() {  
  const ul = document.getElementsByClassName("collections")[0]
  window.go.main.App.GetCollections(false).then(result => {
    if (result === null) {return}; 
    for (var i = 0; i <result.length; i++) {
        var li = document.createElement("li");
        li.setAttribute('class', "collection")
        li.innerHTML = `<div class="collection-id">${result[i]}</div>
        <div class="collection-btns"> 
            <button>details</button>
            <button id="${result[i]}">select</button> 
        </div>`
        ul.appendChild(li)
        document.getElementById(result[i]).onclick = function () {setCollection(this.id)};
    };
    }).catch(err => {
      console.log(err);
    }).finally(() => {
      console.log("finished GetCollections")
    });    
}

function setCollection(id) {
  document.getElementById("active-collection").innerText = id;
}




//    Folders
//   ========================================== 

async function loadFolderHTML() {
  await  fetch('templates/folders.html')
    .then(response=> response.text())
    .then(text=> document.getElementsByClassName('content')[0].innerHTML = text);

  document.getElementById("goback").addEventListener("click", goBack)
}

var FullPath;

async function LoadFolders(id) {
  await loadFolderHTML();
  await window.go.main.App.GetCwd().then(result => {
    FullPath = result;
  }).catch(err => {
    console.log(err);
  }).finally(() => {
    console.log("finished GetCWd")
  });

  document.getElementById("FullPath").innerText = FullPath;
  listFolders("")
  highlight('images')
}

function listFolders(folder) {
  if (folder !== "") {  
    FullPath += "\\" +document.getElementById(folder).innerText;
    document.getElementById("FullPath").innerText = FullPath
  }
  
  window.go.main.App.ListFolders(FullPath).then(result => {
    const ul = document.getElementById("folders");
    ul.innerHTML = "";
    console.log(result )
    if (result === null) {return} 
    for (let i = 0; i <result.length; i++) {
      var li = document.createElement("li");
      li.setAttribute('id', "folder-"+i)
      li.onclick = function() {listFolders(this.id)}
      li.innerText = result[i]
      ul.appendChild(li)
    }
  }).catch(err => {
    console.log(err);
  }).finally(() => {
    console.log("finished! ListFolders")
  });
}



async function goBack() {
  const ul = document.getElementById("folders");
  ul.innerHTML = "";

 await window.go.main.App.GoBack(FullPath).then(result => {
    document.getElementById("FullPath").innerText = result;
    FullPath  = result;
    listFolders("") 
  }).catch(err => {
    console.log(err);
  }).finally(() => {
    console.log("finished GoBack")
  });
}


//    Highlight
//   ========================================== 

function highlight(id) {
  var active = document.getElementsByClassName("active");
  if (active.length !== 0) {
      active[0].className = active[0].className.replace("active", "") 
  }
  document.getElementById(id).className = "active"
}

//    Load Templates
//   ========================================== 
  
document.getElementById('collections').addEventListener('click',async function(){
  await fetch('templates/collections.html')
  .then(response=> response.text())
  .then(text=> document.getElementsByClassName('content')[0].innerHTML = text);
  listCollections()
  highlight('collections')
})
 


// document.getElementById('create').addEventListener('click', function(){
//   fetch('templates/create-collection.html')
//   .then(response=> response.text())
//   .then(text=> document.getElementsByClassName('content')[0].innerHTML = text);
// })
 
// document.getElementById('delete').addEventListener('click', function(){
//   fetch('templates/delete-collection.html')
//   .then(response=> response.text())
//   .then(text=> document.getElementsByClassName('content')[0].innerHTML = text);
// })

// document.getElementById('list').addEventListener('click', async function(){
//   await fetch('templates/list-collections.html')
//   .then(response=> response.text())
//   .then(text=> document.getElementsByClassName('content')[0].innerHTML = text);
//   ListCollections()
// })

document.getElementById('images').addEventListener('click',  function(){
   fetch('templates/images.html')
  .then(response=> response.text())
  .then(text=> document.getElementsByClassName('content')[0].innerHTML = text);
  highlight('images')

})

document.getElementById('matches').addEventListener('click', function(){
  fetch('templates/matches.html')
  .then(response=> response.text())
  .then(text=> document.getElementsByClassName('content')[0].innerHTML = text);
  highlight('matches')

})

document.getElementById('nomatches').addEventListener('click', function(){
  fetch('templates/nomatches.html')
  .then(response=> response.text())
  .then(text=> document.getElementsByClassName('content')[0].innerHTML = text);
  highlight('nomatches')
})


document.getElementById('savematches').addEventListener('click', function(){
  fetch('templates/save-matches.html')
  .then(response=> response.text())
  .then(text=> document.getElementsByClassName('content')[0].innerHTML = text);
})

