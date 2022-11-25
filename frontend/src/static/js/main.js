//  Collections
//  ========================================== 

document.getElementById('collections').addEventListener('click',async function(){
  await fetch('templates/collections.html')
  .then(response=> response.text())
  .then(text=> document.getElementsByClassName('content')[0].innerHTML = text);
  listCollections(false)
  document.getElementById("refresh-collections").addEventListener("click", function() {listCollections(true)})
  document.getElementById("new-collection").addEventListener("click", function() {createCollection()})
  highlight('collections')
})

//refresh Boolean
function listCollections(refresh) {  
  const ul = document.getElementsByClassName("collections")[0]
  ul.innerHTML = ""
  window.go.main.App.GetCollections(refresh).then(result => {
    if (result === null) {
      document.getElementsByClassName('content')[0].innerHTML += "<h2>No collections found</h2>"
      return
    }; 
    for (var i = 0; i <result.length; i++) {
        var li = document.createElement("li");
        li.setAttribute('class', "collection")
        li.innerHTML = `<div class="collection-id">${result[i]}</div>
        <div class="collection-btns"> 
            <button class="collection-details">details</button>
            <button id="${result[i]}">select</button> 
        </div>`
        ul.appendChild(li)
        document.getElementById(result[i]).onclick = function () {setCollection(this.id)};
        document.getElementsByClassName("collection-details")[i].addEventListener("click", function () {alert("TODO :)")});
      };
    }).catch(err => {
      alert(err);
    }).finally(() => {
      console.log("finished GetCollections")
    });    
}


function setCollection(id) {
  document.getElementById("active-collection").innerText = id;
}

async function createCollection() {
  let collection = prompt("New collection ID:");
  if (collection == null || collection == "") {
    alert("operation cancelled")
    return
  } 
  await window.go.main.App.CreateCollection(collection).catch(err => {
      alert(err);
    }).finally(() => {
      console.log("finished GetCollections")
    });    
    listCollections(true)
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

//   Images
//   ========================================== 

document.getElementById('images').addEventListener('click',  async function(){
  await fetch('templates/images.html')
 .then(response=> response.text())
 .then(text=> document.getElementsByClassName('content')[0].innerHTML = text);
 highlight('images')
 document.getElementById("data-dir").addEventListener("click", searchFaces)
})

async function searchFaces() {
  let collection = document.getElementById("active-collection").innerText
  await window.go.main.App.SearchFaces(collection).catch(err => {
    alert(err);
  }).finally(() => {
    alert("succefully uploaded images")
  });
}

//    Results 
//   ========================================== 

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

