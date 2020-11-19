'use strict';
let globalData;
let ws;


$(document).ready(function() {

    let url = 'ws://' + location.host + location.pathname + location.search + 'test/ws'
    ws = new WebSocket(url)

    ws.onopen = function() {
        console.log('connected')
    };

    ws.onclose = function(evt) {
        console.log('disconnected', evt);
    };

    ws.onmessage = function(evt) {
        let allData = JSON.parse(evt.data);
        let data = allData.data;
        console.log(allData);

        switch (allData.type) {
            case 'custInfo':
                {
                    let customers = data.custInfo.sort(sortByName);
                    fillTalbeCustomer(customers, true);
                    break;
                }
            case 'custUpdate':
                {
                    let customers = data.custInfo.sort(sortByName);
                    fillTalbeCustomer(customers, false);
                    break;
                }
            case 'error':
                {
                    $('#infoAlert').text(data.message.error)
                    $('#infoAlert').fadeTo(3000, 500).slideUp(500, function() {
                        $('#infoAlert').slideUp(500);
                    });
                }

            default:
                break;
        };
    };


    //кнопка Создать
    $('#bt_create').on('click', function() {
        setCreateDialog();
        $('#custDialog').dialog('open');
        $('#custDialog').dialog({
            buttons: {
                "Отправить": function() {
                    sendCustomerDialog("createCustomer");
                },
            },
        });
    });

    //кнопка Обновить
    $('#bt_update').prop('disabled', true);
    $('#bt_update').on('click', function() {
        setUpdateDialog();
        $('#custDialog').dialog('open');
        $('#custDialog').dialog({
            buttons: {
                "Отправить": function() {
                    sendCustomerDialog("updateCustomer")
                },
            },
        });
    });

    //кнопка Удалить
    $('#bt_delete').prop('disabled', true);
    $('#bt_delete').on('click', function() {
        deleteB();
    });
});

function sortByName(a, b) {
    return a.name - b.name;
};

//fillTalbeCustomer заполнить таблицу клиентов
function fillTalbeCustomer(customers, firstFlag) {
    let $table = $('#table');
    let selected = $table.bootstrapTable('getSelections');
    let toWrite = [];
    customers.forEach(cust => {
        let temp = {
            id: cust.id,
            check: false,
            name: cust.name,
            address: cust.address,
            numS: cust.servers.length,
            phone: cust.phone,
            email: cust.email,
        };
        if (selected.length === 1) {
            if (cust.id === selected[0].id) {
                temp.check = true;
            }
        }
        toWrite.push(temp);
    });
    $table.bootstrapTable('load', toWrite);
    $table.bootstrapTable('hideColumn', 'id');
    $table.bootstrapTable('scrollTo', 'top');

    $table.unbind().on('click', function() {
        $('#bt_update').prop('disabled', false);
        $('#bt_delete').prop('disabled', false);
    });
};

//setCreateDialog диалог при создании клиента
function setCreateDialog() {
    $('#custDialog').dialog({
        autoOpen: false,
        resizable: false,
    });
    $('#name').val("");
    $('#address').val("");
    $('#phone').val("");
    $('#email').val("");
};

//setUpdateDialog диалог при обновлении клиента
function setUpdateDialog() {
    $('#custDialog').dialog({
        autoOpen: false,
        resizable: false,
    });
    let $table = $('#table');
    let selected = $table.bootstrapTable('getSelections');
    $('#name').val(selected[0].name);
    $('#address').val(selected[0].address);
    $('#phone').val(selected[0].phone);
    $('#email').val(selected[0].email);
}

//sendCustomerDialog отпрака информации из диалога клиентов на сервер
function sendCustomerDialog(typeD) {
    let cForm = $('#custForm')
    if (!cForm[0].checkValidity()) {
        cForm[0].classList.add('was-validated');
        return
    }
    let selected = $('#table').bootstrapTable('getSelections');
    let toSend = {
        type: typeD,
        name: $('#name').val(),
        address: $('#address').val(),
        phone: $('#phone').val(),
        email: $('#email').val(),
    };
    if (selected.length === 0) {
        toSend.id = 0;
    } else {
        toSend.id = selected[0].id;
    }
    ws.send(JSON.stringify(toSend));
    $('#custDialog').dialog('close');
};

//deleteB удаление клиента
function deleteB() {
    let selected = $('#table').bootstrapTable('getSelections');
    let toSend = {
        type: "deleteCustomer",
        id: selected[0].id,
    };
    $('#bt_update').prop('disabled', true);
    $('#bt_delete').prop('disabled', true);
    ws.send(JSON.stringify(toSend));
};