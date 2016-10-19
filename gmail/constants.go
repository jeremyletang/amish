package gmail

const authSuccess = `
<style type="text/css">
  @import url(https://fonts.googleapis.com/css?family=Roboto:400,500);
  /*@import url(https://piemapping.com/fonts/DINRoundPro-Medi.otf?family=DINRoundPro-Medi);*/

  html {
    box-sizing: border-box;
  }
  *, *:before, *:after {
    box-sizing: inherit;
  }

  body,
  html,
  .body {
    background: #F7F8FA !important;
    font-family: Roboto, Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    color: #434b5a;
  }

  h1 {
    font-weight: 500;
    font-size: 23px;
    margin-top: 0;
  }

  .separator {
    width: 40px;
    border-bottom: 2px solid #4A90E2;
    margin: 25px auto;
    border-radius: 1px;
  }

  p {
    line-height: 24px;
    margin: 0;
  }

  .container-img {
    margin: 30px 0;
  }

  .card {
    background-color: #fff;
    padding: 50px;
    border-radius: 5px;
    max-width: 600px;
    margin: 0 auto;
    border: 1px solid #D4D9DD;
  }

  a {
    text-decoration: none;
  }

  a.btn {
    display: inline-block;
  }

  .btn {
    border-radius: 5px;
    border: none;
    line-height: 48px;
    padding: 0 20px;
    font-size: 14px;
    font-weight: 500;
  }

  .btn.blue {
    color: #fff;
    background-color: #4A90E2;
    padding: 0 68px;
    max-width: 100%;
    border: solid 1px #2d62a0;
  }

  .footer {
    font-size: 14px;
    color: #989CA5;
    margin-top: 25px;
    padding: 0;
  }

  .footer > item:before {
    content: "\25CF";
    padding-left: 3px;
    padding-right: 3px;
    font-size: 6px;
    vertical-align: middle;
  }

  .footer > item.first-item:before {
    content: "";
    padding-left: 0;
    padding-right: 0;
    font-size: 6px;
    vertical-align: middle;
  }

  @media screen and (max-width: 399px) {
    .card {
      padding: 30px 20px;
    }

    .btn.blue {
      padding: 0 20px;
    }
  }
</style>

<container class="body-border">
  <row>
    <columns>

      <spacer size="32"></spacer>

      <center class="container-img">
          <h1>Authentication success</h1>
      </center>

      <spacer size="16"></spacer>

      <center class="card">
        <h1>
          You Successfully authentified with AMISH.
        </h1>

        <div class="separator"></div>

       </center>
    </columns>

    <center>
      <!--
      <menu class="text-center footer">
        <item class="first-item" href="#">Pie Mapping</item>
        <item href="#">TechHub London</item>
        <item href="#">20 Ropemaker St</item>
        <item href="#">Moorgate EC2Y 9AR</item>
      </menu>
      -->
    </center>

  </row>

  <spacer size="16"></spacer>
</container>
`
