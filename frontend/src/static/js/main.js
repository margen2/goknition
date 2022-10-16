import { ListCollections } from './collections.js'
// //Global JS function for greeting
// function greet() {
//   //Get user input
//   let inputName = document.getElementById("name").value;

//   //Call Go Greet function
//   window.go.main.App.Greet(inputName).then(result => {
//     //Display result from Go
//     document.getElementById("result").innerHTML = result;
//   }).catch(err => {
//     console.log(err);
//   }).finally(() => {
//     console.log("finished!")
//   });
// }

document.getElementById('create').addEventListener('click', function(){
  fetch('templates/create-collection.html')
  .then(response=> response.text())
  .then(text=> document.getElementsByTagName('main')[0].innerHTML = text);
})
 
document.getElementById('delete').addEventListener('click', function(){
  fetch('templates/delete-collection.html')
  .then(response=> response.text())
  .then(text=> document.getElementsByTagName('main')[0].innerHTML = text);
})

document.getElementById('list').addEventListener('click', async function(){
  await fetch('templates/list-collections.html')
  .then(response=> response.text())
  .then(text=> document.getElementsByTagName('main')[0].innerHTML = text);
  ListCollections()
})

document.getElementById('images').addEventListener('click', function(){
  fetch('templates/images.html')
  .then(response=> response.text())
  .then(text=> document.getElementsByTagName('main')[0].innerHTML = text);
})

document.getElementById('matches').addEventListener('click', function(){
  fetch('templates/matches.html')
  .then(response=> response.text())
  .then(text=> document.getElementsByTagName('main')[0].innerHTML = text);
})

document.getElementById('nomatches').addEventListener('click', function(){
  fetch('templates/nomatches.html')
  .then(response=> response.text())
  .then(text=> document.getElementsByTagName('main')[0].innerHTML = text);
})


document.getElementById('savematches').addEventListener('click', function(){
  fetch('templates/save-matches.html')
  .then(response=> response.text())
  .then(text=> document.getElementsByTagName('main')[0].innerHTML = text);
})

