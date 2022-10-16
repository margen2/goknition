
export function ListCollections() {  
    var ul = document.getElementById("collections")
    window.go.main.App.GetCollections(true).then(result => {
        //Display result from Go
        for (let i = 0; i <result.length; i++) {
            var li = document.createElement("li")
            li.setAttribute('id', "collection-"+i)
            li.innerText = result[i]
            ul.appendChild(li)
        }
      }).catch(err => {
        console.log(err);
      }).finally(() => {
        console.log("finished!")
      });    
}
function greet() {
    //Get user input
    let inputName = document.getElementById("name").value;
  
    //Call Go Greet function
    window.go.main.App.Greet(inputName).then(result => {
      //Display result from Go
      document.getElementById("result").innerHTML = result;
    }).catch(err => {
      console.log(err);
    }).finally(() => {
      console.log("finished!")
    });
  }