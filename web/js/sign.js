async function signMsg(msgParams, from, state) {
    web3.currentProvider.sendAsync({
      method: 'eth_signTypedData',
      params: [msgParams, from],
      from: from,
    }, async function (err, result) {
      if (err) return console.error(err)
      if (result.error) {
        return console.error(result.error.message)
      }


      await fetch("http://127.0.0.1/v1/session/check_signature", {
        method: "post",
        body: JSON.stringify({
          signature: result.result,
        }),
        headers: {
          "Authorization": state,
          "content-type": "application/json"
        },
        credentials: 'include'    
      });

      window.open('','_self').close()
    })
}