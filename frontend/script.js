var input = document.getElementById("input-uid")
var button = document.getElementById("submit-uid-btn")


 button.onclick = function() {

    console.log(input.value)

       fetch('http://localhost:8081/order/' + input.value, {
          mode: 'cors',
          headers: {
            'Origin': 'http://localhost:8081'
          }
       })
    .then(result => result.json())
    .then((orderData) => {
      console.log(orderData)

          htmlOrderString = `
      <div class="order-section">
        <h2>Order Information</h2>
        <p><div>Order ID: ${orderData.order_uid}</div></p>
        <p><div>Track Number: ${orderData.track_number}</div></p>
        <p><div>Date Created: ${orderData.date_created}</div></p>
      </div>

      <div class="order-section">
        <h2>Delivery Details</h2>
        <p><div>Name: ${orderData.delivery.name}</div></p>
        <p><div>Address: ${orderData.delivery.address}, ${orderData.delivery.city}</div></p>
        <p><div>Phone: ${orderData.delivery.phone}</div></p>
        <p><div>Email: ${orderData.delivery.email}</div></p>
      </div>

      <div class="order-section">
        <h2>Payment Information</h2>
        <p><div>Amount: ${orderData.payment.amount} ${orderData.payment.currency}</div></p>
        <p><div>Payment Date: ${orderData.payment.payment_dt}</div></p>
        <p><div>Provider: ${orderData.payment.provider}</div></p>
      </div>
    `;

    htmlItemString = ""

      orderData.items.forEach(item => {
        htmlItemString += `
      <div class="item-card">
        <h2>${item.name} (${item.brand})</h2>
        <p><div>Price: ${item.price} (${item.sale}% off)</div> </p>
        <p><div>Total: ${item.total_price}</div></p>
      </div>
    `;
        
      });
      
    container.innerHTML = htmlOrderString;
    item_container.innerHTML = htmlItemString;

    })
    .catch(err => console.error(err));

    var container = document.getElementById("order-container");
    var item_container = document.getElementById("item-container");
    };
