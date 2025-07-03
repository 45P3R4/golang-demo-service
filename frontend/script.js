var input = document.getElementById("input-uid")
var button = document.getElementById("submit-uid-btn")


 button.onclick = function() {

    console.log(input.value)

       fetch('http://localhost:8081/order/' + input.value)
    .then(result => result.json())
    .then((orderData) => {
      console.log(orderData)

          htmlOrderString = `
      <div class="order-section">
        <h2>Order Information</h2>
        <p><span class="label">Order ID:</span> ${orderData.order_uid}</p>
        <p><span class="label">Track Number:</span> ${orderData.track_number}</p>
        <p><span class="label">Date Created:</span> ${orderData.date_created}</p>
      </div>

      <div class="order-section">
        <h2>Delivery Details</h2>
        <p><span class="label">Name:</span> ${orderData.delivery.name}</p>
        <p><span class="label">Address:</span> ${orderData.delivery.address}, ${orderData.delivery.city}</p>
        <p><span class="label">Phone:</span> ${orderData.delivery.phone}</p>
        <p><span class="label">Email:</span> ${orderData.delivery.email}</p>
      </div>

      <div class="order-section">
        <h2>Payment Information</h2>
        <p><span class="label">Amount:</span> ${orderData.payment.amount} ${orderData.payment.currency}</p>
        <p><span class="label">Payment Date:</span> ${orderData.payment.payment_dt}</p>
        <p><span class="label">Provider:</span> ${orderData.payment.provider}</p>
      </div>
    `;

    //   htmlItemString = `
    //   <div class="order-item">
    //     <p><span class="label">Product:</span> ${item.name} (${item.brand})</p>
    //     <p><span class="label">Price:</span> ${item.price} (${item.sale}% off)</p>
    //     <p><span class="label">Total:</span> ${item.total_price}</p>
    //   </div>
    // `;
    container.innerHTML = htmlOrderString

    })
    .catch(err => console.error(err));

    var container = document.getElementById("order-container");
    };



    




