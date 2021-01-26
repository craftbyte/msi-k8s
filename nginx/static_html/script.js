async function initialize() {
    document.getElementById("addBtn").addEventListener("click", addTask);
    document.getElementById("addForm").addEventListener("submit", addTask);
    let tasks = await (await fetch('/api/tasks', {
        cache: 'no-cache'
    })).json()
    for (let task of tasks) {
        showTask(task)
    }
}

function showTask(task) {
    let taskList = document.getElementById("tasks")
    let li = document.createElement("li");
    li.addEventListener("dblclick", removeTask);
    let elem = document.createElement("label");
    let check = document.createElement("input");
    check.type = "checkbox";
    check.checked = task.done;
    check.dataset.id = task.id;
    check.addEventListener("change", updateTask)
    elem.appendChild(check);
    let span = document.createElement("span");
    span.innerText = task.name;
    elem.appendChild(span);
    li.appendChild(elem);
    taskList.appendChild(li);
}

async function addTask(event) {
    event.preventDefault();
    let name = document.getElementById("add").value;
    document.getElementById("add").value = "";
    let task = await(await fetch('/api/tasks', {
        method: 'POST',
        cache: 'no-cache',
        body: JSON.stringify({name})
    })).json()
    showTask(task);
}

async function updateTask() {
    let id = this.dataset.id;
    let checked = this.checked;
    let task = await(await fetch('/api/tasks/'+id, {
        method: 'PATCH',
        cache: 'no-cache',
        body: JSON.stringify({done: checked})
    })).json()
    this.checked = task.done;
}

async function removeTask(event) {
    event.preventDefault();
    let id = this.querySelector('input[type="checkbox"]').dataset.id;
    await fetch('/api/tasks/'+id, {
        method: 'DELETE',
        cache: 'no-cache'
    })
    this.parentNode.removeChild(this);
}
document.addEventListener("DOMContentLoaded", initialize)