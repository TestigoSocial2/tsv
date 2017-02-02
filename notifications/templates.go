package notifications

import (
	"bytes"
	"html/template"
)

// PrepareContent will return a notification's content using the provided string and data
func PrepareContent(tpl *template.Template, data map[string]interface{}) (string, error) {
	buf := new(bytes.Buffer)
	err := tpl.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// TSVEmailTemplate defines the default TSV style for email notifications
var TSVEmailTemplate, _ = template.New("tsv").Parse(`
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>Testigo Social 2.0</title>
    <style type="text/css">
      /* CLIENT-SPECIFIC STYLES */
      #outlook a{
        padding:0;
      }
      .ReadMsgBody{
        width:100%;
      }

      .ExternalClass{
        width:100%;
      }

      .ExternalClass, .ExternalClass p, .ExternalClass span, .ExternalClass font, .ExternalClass td, .ExternalClass div {
        line-height: 100%;
      }

      body, table, td, p, a, li, blockquote{
        -webkit-text-size-adjust:100%; -ms-text-size-adjust:100%;
      }

      table, td{
        mso-table-lspace:0pt; mso-table-rspace:0pt;
      }

      img{
        -ms-interpolation-mode:bicubic;
      }

      /* RESET STYLES */
      body{margin:0; padding:0;}
      img{border:0; height:auto; line-height:100%; outline:none; text-decoration:none;}
      table{border-collapse:collapse !important;}
      body, #bodyTable, #bodyCell{height:100% !important; margin:0; padding:0; width:100% !important;}

      /* TEMPLATE STYLES */
      #bodyCell{padding:20px;}
      #templateContainer{width:600px;}

      body, #bodyTable{
        background-color:#DEE0E2;
      }

      #bodyCell{
        border-top:4px solid #BBBBBB;
      }

      #templateContainer{
        border:1px solid #BBBBBB;
      }

      h1{
        color:#202020 !important;
        display:block;
        font-family:Helvetica;
        font-size:26px;
        font-style:normal;
        font-weight:bold;
        line-height:100%;
        letter-spacing:normal;
        margin-top:0;
        margin-right:0;
        margin-bottom:10px;
        margin-left:0;
        text-align:left;
      }

      h2{
        color:#404040 !important;
        display:block;
        font-family:Helvetica;
        font-size:20px;
        font-style:normal;
        font-weight:bold;
        line-height:100%;
        letter-spacing:normal;
        margin-top:0;
        margin-right:0;
        margin-bottom:10px;
        margin-left:0;
        text-align:left;
      }

      h3{
        color:#606060 !important;
        display:block;
        font-family:Helvetica;
        font-size:16px;
        font-style:italic;
        font-weight:normal;
        line-height:100%;
        letter-spacing:normal;
        margin-top:0;
        margin-right:0;
        margin-bottom:10px;
        margin-left:0;
        text-align:left;
      }

      h4{
        color:#808080 !important;
        display:block;
        font-family:Helvetica;
        font-size:14px;
        font-style:italic;
        font-weight:normal;
        line-height:100%;
        letter-spacing:normal;
        margin-top:0;
        margin-right:0;
        margin-bottom:10px;
        margin-left:0;
        text-align:left;
      }

      #templatePreheader{
        background-color:#F4F4F4;
        border-bottom:1px solid #CCCCCC;
      }

      .preheaderContent{
        color:#808080;
        font-family:Helvetica;
        font-size:10px;
        line-height:125%;
        text-align:left;
      }

      /* Yahoo! Mail Override */
      .preheaderContent a:link, .preheaderContent a:visited, .preheaderContent a .yshortcuts {
        color:#606060;
        font-weight:normal;
        text-decoration:underline;
      }

      #templateHeader{
        background-color:#F4F4F4;
        border-top:1px solid #FFFFFF;
        border-bottom:1px solid #CCCCCC;
      }

      .headerContent{
        color:#505050;
        font-family:Helvetica;
        font-size:20px;
        font-weight:bold;
        line-height:100%;
        padding-top:0;
        padding-right:0;
        padding-bottom:0;
        padding-left:0;
        text-align:left;
        vertical-align:middle;
      }

      /* Yahoo! Mail Override */
      .headerContent a:link, .headerContent a:visited, .headerContent a .yshortcuts {
        color:#EB4102;
        font-weight:normal;
        text-decoration:underline;
      }

      #headerImage{
        height:auto;
        max-width:600px;
      }

      #templateBody{
        background-color:#F4F4F4;
        border-top:1px solid #FFFFFF;
        border-bottom:1px solid #CCCCCC;
      }

      .bodyContent{
        color:#505050;
        font-family:Helvetica;
        font-size:16px;
        line-height:150%;
        padding-top:20px;
        padding-right:20px;
        padding-bottom:20px;
        padding-left:20px;
        text-align:left;
      }

      /* Yahoo! Mail Override */
      .bodyContent a:link, .bodyContent a:visited, .bodyContent a .yshortcuts {
        color:#EB4102;
        font-weight:normal;
        text-decoration:underline;
      }

      .bodyContent img{
        display:inline;
        height:auto;
        max-width:560px;
      }

      .templateColumnContainer{width:260px;}

      #templateColumns{
        background-color:#F4F4F4;
        border-top:1px solid #FFFFFF;
        border-bottom:1px solid #CCCCCC;
      }

      .leftColumnContent{
        color:#505050;
        font-family:Helvetica;
        font-size:14px;
        line-height:150%;
        padding-top:0;
        padding-right:20px;
        padding-bottom:20px;
        padding-left:20px;
        text-align:left;
      }

      /* Yahoo! Mail Override */
      .leftColumnContent a:link, .leftColumnContent a:visited, .leftColumnContent a .yshortcuts {
        color:#EB4102;
        font-weight:normal;
        text-decoration:underline;
      }

      .rightColumnContent{
        color:#505050;
        font-family:Helvetica;
        font-size:14px;
        line-height:150%;
        padding-top:0;
        padding-right:20px;
        padding-bottom:20px;
        padding-left:20px;
        text-align:left;
      }

      /* Yahoo! Mail Override */
      .rightColumnContent a:link, .rightColumnContent a:visited, .rightColumnContent a .yshortcuts {
        color:#EB4102;
        font-weight:normal;
        text-decoration:underline;
      }

      .leftColumnContent img, .rightColumnContent img{
        display:inline;
        height:auto;
        max-width:260px;
      }

      #templateFooter{
        background-color:#222222;
        border-top:1px solid #FFFFFF;
      }

      .footerContent{
        color:#808080;
        font-family:Helvetica;
        font-size:10px;
        line-height:150%;
        padding-top:20px;
        padding-right:20px;
        padding-bottom:20px;
        padding-left:20px;
        text-align:left;
      }

      /* Yahoo! Mail Override */
      .footerContent a:link, .footerContent a:visited, .footerContent a .yshortcuts, .footerContent a span {
        color:#606060;
        font-weight:normal;
        text-decoration:underline;
      }

      @media only screen and (max-width: 480px){
        body, table, td, p, a, li, blockquote{-webkit-text-size-adjust:none !important;}
        body{width:100% !important; min-width:100% !important;}
        #bodyCell{padding:10px !important;}
        #templateContainer{
          max-width:600px !important;
          width:100% !important;
        }

        h1{
          font-size:24px !important;
          line-height:100% !important;
        }

        h2{
          font-size:20px !important;
          line-height:100% !important;
        }

        h3{
          font-size:18px !important;
          line-height:100% !important;
        }

        h4{
          font-size:16px !important;
          line-height:100% !important;
        }

        #templatePreheader{display:none !important;}

        #headerImage{
          height:auto !important;
          max-width:600px !important;
          width:100% !important;
        }

        .headerContent{
          font-size:20px !important;
          line-height:125% !important;
        }

        .bodyContent{
          font-size:18px !important;
          line-height:125% !important;
        }

        .templateColumnContainer{display:block !important; width:100% !important;}

        .columnImage{
          height:auto !important;
          max-width:480px !important;
          width:100% !important;
        }

        .leftColumnContent{
          font-size:16px !important;
          line-height:125% !important;
        }

        .rightColumnContent{
          font-size:16px !important;
          line-height:125% !important;
        }

        .footerContent{
          font-size:14px !important;
          line-height:115% !important;
        }

        .footerContent a{display:block !important;}
      }
    </style>
  </head>
  <body leftmargin="0" marginwidth="0" topmargin="0" marginheight="0" offset="0">
    <center>
      <table align="center" border="0" cellpadding="0" cellspacing="0" height="100%" width="100%" id="bodyTable">
        <tr>
          <td align="center" valign="top" id="bodyCell">
            <table border="0" cellpadding="0" cellspacing="0" id="templateContainer">
              <tr>
                <td align="center" valign="top"></td>
              </tr>
              <tr>
                <td align="center" valign="top">
                  <table border="0" cellpadding="0" cellspacing="0" width="100%" id="templateHeader">
                    <tr>
                      <td valign="top" class="headerContent">
                        <img src="https://s3.amazonaws.com/tm.org.mx/tsv/email_header_01.png" style="max-width:600px;" id="headerImage" mc:label="header_image" mc:edit="header_image" mc:allowdesigner mc:allowtext />
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
              <tr>
                <td align="center" valign="top">
                  <table border="0" cellpadding="0" cellspacing="0" width="100%" id="templateBody">
                    <tr>
                      <td valign="top" class="bodyContent" mc:edit="body_content">
                        <h1>{{.Title}}</h1>
                        <h3>{{.Subtitle}}</h3>
                        {{.Content}}
                        <br />
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
              <!--
              <tr>
                <td align="center" valign="top">
                  <table border="0" cellpadding="0" cellspacing="0" width="100%" id="templateColumns">
                    <tr mc:repeatable>
                      <td align="center" valign="top" class="templateColumnContainer" style="padding-top:20px;">
                        <table border="0" cellpadding="20" cellspacing="0" width="100%">
                          <tr>
                            <td class="leftColumnContent">
                              <img src="http://placehold.it/260x120" style="max-width:260px;" class="columnImage" mc:label="left_column_image" mc:edit="left_column_image" />
                            </td>
                          </tr>
                          <tr>
                            <td valign="top" class="leftColumnContent" mc:edit="left_column_content">
                              <h3>First Column</h3>
                              Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
                            </td>
                          </tr>
                        </table>
                      </td>
                      <td align="center" valign="top" class="templateColumnContainer" style="padding-top:20px;">
                        <table border="0" cellpadding="20" cellspacing="0" width="100%">
                          <tr>
                            <td class="rightColumnContent">
                              <img src="http://placehold.it/260x120" style="max-width:260px;" class="columnImage" mc:label="right_column_image" mc:edit="right_column_image" />
                            </td>
                          </tr>
                          <tr>
                            <td valign="top" class="rightColumnContent" mc:edit="right_column_content">
                              <h3>Second Column</h3>
                              Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
                            </td>
                          </tr>
                        </table>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
              -->
              <tr>
                <td align="center" valign="top">
                  <table border="0" cellpadding="0" cellspacing="0" width="100%" id="templateFooter">
                    <tr>
                      <td valign="top" class="footerContent" mc:edit="footer_content00">
                        <a href="https://www.twitter.com/testigo_social">Twitter</a>&nbsp;&nbsp;&nbsp;<a href="https://www.facebook.com/testigosocial">Facebook</a>
                      </td>
                    </tr>
                    <tr>
                      <td valign="top" class="footerContent" style="padding-top:0;" mc:edit="footer_content01">
                        <strong>Transparencia Mexicana AC &copy;</strong>
                        <br />
                        Dulce Olivia 73<br />
                        Colonia Villa Coyoacán<br />
                        Delegación Coyoacán<br />
                        CP 04000<br />
                        México, Ciudad de México
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>
          </td>
        </tr>
      </table>
    </center>
  </body>
</html>
`)

// SimpleEmailTemplate simple no-images template with sidebar content
var SimpleEmailTemplate, _ = template.New("simple").Parse(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>Testigo Social 2.0</title>
    <style type="text/css">
    /* Client-specific Styles */
    #outlook a{padding:0;} /* Force Outlook to provide a "view in browser" button. */
    body{width:100% !important;} .ReadMsgBody{width:100%;} .ExternalClass{width:100%;} /* Force Hotmail to display emails at full width */
    body{-webkit-text-size-adjust:none;} /* Prevent Webkit platforms from changing default text sizes. */

    /* Reset Styles */
    body{margin:0; padding:0;}
    img{border:0; height:auto; line-height:100%; outline:none; text-decoration:none;}
    table td{border-collapse:collapse;}
    #backgroundTable{height:100% !important; margin:0; padding:0; width:100% !important;}

    body, #backgroundTable{
      background-color:#FAFAFA;
    }

    #templateContainer{
      border:0;
    }

    h1, .h1{
      color:#202020;
      display:block;
      font-family:Arial;
      font-size:40px;
      font-weight:bold;
      line-height:100%;
      margin-top:2%;
      margin-right:0;
      margin-bottom:1%;
      margin-left:0;
      text-align:left;
    }

    h2, .h2{
      color:#404040;
      display:block;
      font-family:Arial;
      font-size:18px;
      font-weight:bold;
      line-height:100%;
      margin-top:2%;
      margin-right:0;
      margin-bottom:1%;
      margin-left:0;
      text-align:left;
    }

    h3, .h3{
      color:#606060;
      display:block;
      font-family:Arial;
      font-size:16px;
      font-weight:bold;
      line-height:100%;
      margin-top:2%;
      margin-right:0;
      margin-bottom:1%;
      margin-left:0;
      text-align:left;
    }

    h4, .h4{
      color:#808080;
      display:block;
      font-family:Arial;
      font-size:14px;
      font-weight:bold;
      line-height:100%;
      margin-top:2%;
      margin-right:0;
      margin-bottom:1%;
      margin-left:0;
      text-align:left;
    }

    #templatePreheader{
      background-color:#FAFAFA;
    }

    .preheaderContent div{
      color:#707070;
      font-family:Arial;
      font-size:10px;
      line-height:100%;
      text-align:left;
    }

    /* Yahoo! Mail Override */
    .preheaderContent div a:link, .preheaderContent div a:visited, .preheaderContent div a .yshortcuts {
      color:#336699;
      font-weight:normal;
      text-decoration:underline;
    }

    #social div{
      text-align:right;
    }

    #templateHeader{
      background-color:#FFFFFF;
      border-bottom:5px solid #505050;
    }

    .headerContent{
      color:#202020;
      font-family:Arial;
      font-size:34px;
      font-weight:bold;
      line-height:100%;
      padding:10px;
      text-align:right;
      vertical-align:middle;
    }

    /* Yahoo! Mail Override */
    .headerContent a:link, .headerContent a:visited, .headerContent a .yshortcuts {
      color:#336699;
      font-weight:normal;
      text-decoration:underline;
    }

    #headerImage{
      height:auto;
      max-width:600px !important;
    }

    #templateContainer, .bodyContent{
      background-color:#FDFDFD;
    }

    .bodyContent div{
      color:#505050;
      font-family:Arial;
      font-size:14px;
      line-height:150%;
      text-align:justify;
    }

    /* Yahoo! Mail Override */
    .bodyContent div a:link, .bodyContent div a:visited, .bodyContent div a .yshortcuts {
      color:#336699;
      font-weight:normal;
      text-decoration:underline;
    }

    .bodyContent img{
      display:inline;
      height:auto;
    }

    #templateSidebar{
      background-color:#FDFDFD;
    }

    .sidebarContent{
      border-left:1px solid #DDDDDD;
    }

    .sidebarContent div{
      color:#505050;
      font-family:Arial;
      font-size:10px;
      line-height:150%;
      text-align:left;
    }

    /* Yahoo! Mail Override */
    .sidebarContent div a:link, .sidebarContent div a:visited, .sidebarContent div a .yshortcuts {
      color:#336699;
      font-weight:normal;
      text-decoration:underline;
    }

    .sidebarContent img{
      display:inline;
      height:auto;
    }

    #templateFooter{
      background-color:#FAFAFA;
      border-top:3px solid #909090;
    }

    .footerContent div{
      color:#707070;
      font-family:Arial;
      font-size:11px;
      line-height:125%;
      text-align:left;
    }

    /* Yahoo! Mail Override */
    .footerContent div a:link, .footerContent div a:visited, .footerContent div a .yshortcuts {
      color:#336699;
      font-weight:normal;
      text-decoration:underline;
    }

    .footerContent img{
      display:inline;
    }

    #social{
      background-color:#FFFFFF;
      border:0;
    }

    #social div{
      text-align:left;
    }

    #utility{
      background-color:#FAFAFA;
      border-top:0;
    }

    #utility div{
      text-align:left;
    }

    #monkeyRewards img{
      max-width:170px !important;
    }
    </style>
  </head>
  <body leftmargin="0" marginwidth="0" topmargin="0" marginheight="0" offset="0">
    <center>
      <table border="0" cellpadding="0" cellspacing="0" height="100%" width="100%" id="backgroundTable">
        <tr>
          <td align="center" valign="top">
            <table border="0" cellpadding="0" cellspacing="0" width="600" id="templateContainer">
              <tr>
                <td align="center" valign="top">
                  <table border="0" cellpadding="0" cellspacing="0" width="600" id="templateHeader">
                    <tr>
                      <td class="headerContent" width="100%" style="padding-left:20px; padding-right:10px;">
                        <div mc:edit="Header_content">
                          <h1>Testigo Social 2.0</h1>
                        </div>
                      </td>
                      <td class="headerContent"></td>
                    </tr>
                  </table>
                </td>
              </tr>
              <tr>
                <td align="center" valign="top">
                  <table border="0" cellpadding="10" cellspacing="0" width="600" id="templateBody">
                    <tr>
                      <td valign="top" class="bodyContent">
                        <table border="0" cellpadding="10" cellspacing="0" width="100%">
                          <tr>
                            <td valign="top" style="padding-right:0;">
                              <div mc:edit="std_content00">
                                <h2 class="h2">{{.Title}}</h2>
                                {{.Content}}
                                <br />
                              </div>
                            </td>
                          </tr>
                        </table>
                      </td>
                      <td valign="top" width="180" id="templateSidebar">
                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                          <tr>
                            <td valign="top">
                              <table border="0" cellpadding="20" cellspacing="0" width="100%" class="sidebarContent">
                                <tr>
                                  <td valign="top" style="padding-right:10px;">
                                    <div mc:edit="std_content01">
                                      <strong>{{.SidebarTitle}}</strong>
                                      <br />
                                      {{.SidebarContent}}
                                    </div>
                                  </td>
                                </tr>
                              </table>
                            </td>
                          </tr>
                        </table>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
              <tr>
                <td align="center" valign="top">
                  <table border="0" cellpadding="0" cellspacing="0" width="600" id="templateFooter">
                    <tr>
                      <td valign="top" class="footerContent">
                        <table border="0" cellpadding="10" cellspacing="0" width="100%">
                          <tr>
                            <td colspan="2" valign="middle" id="social">
                              <div mc:edit="std_social">
                                &nbsp;<a href="https://www.twitter.com/testigo_social">Twitter</a> | <a href="https://www.facebook.com/testigosocial">Facebook</a>
                              </div>
                            </td>
                          </tr>
                          <tr>
                            <td valign="top" width="350">
                              <div mc:edit="std_footer">
                                <strong>Transparencia Mexicana AC</strong>
                                <br />
                                Dulce Olivia 73<br />
                                Colonia Villa Coyoacán<br />
                                Delegación Coyoacán<br />
                                CP 04000<br />
                                México, Ciudad de México
                              </div>
                            </td>
                          </tr>
                        </table>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>
            <br />
          </td>
        </tr>
      </table>
    </center>
  </body>
</html>
`)

// TableEmailTemplate support a slice of items to be rendered as a table
var TableEmailTemplate, _ = template.New("table").Parse(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <meta property="og:title" content="Testigo Social 2.0" />
    <title>Testigo Social 2.0</title>
    <style type="text/css">
    /* Client-specific Styles */
    #outlook a{padding:0;} /* Force Outlook to provide a "view in browser" button. */
    body{width:100% !important;} .ReadMsgBody{width:100%;} .ExternalClass{width:100%;} /* Force Hotmail to display emails at full width */
    body{-webkit-text-size-adjust:none;} /* Prevent Webkit platforms from changing default text sizes. */

    /* Reset Styles */
    body{margin:0; padding:0;}
    img{border:0; height:auto; line-height:100%; outline:none; text-decoration:none;}
    table td{border-collapse:collapse;}
    #backgroundTable{height:100% !important; margin:0; padding:0; width:100% !important;}

    body, #backgroundTable{
      background-color:#FAFAFA;
    }

    #templateContainer{
      border: 1px solid #DDDDDD;
    }

    h1, .h1{
      color:#202020;
      display:block;
      font-family:Arial;
      font-size:34px;
      font-weight:bold;
      line-height:100%;
      margin-top:0;
      margin-right:0;
      margin-bottom:10px;
      margin-left:0;
      text-align:left;
    }

    h2, .h2{
      color:#202020;
      display:block;
      font-family:Arial;
      font-size:30px;
      font-weight:bold;
      line-height:100%;
      margin-top:0;
      margin-right:0;
      margin-bottom:10px;
      margin-left:0;
      text-align:left;
    }

    h3, .h3{
      color:#202020;
      display:block;
      font-family:Arial;
      font-size:26px;
      font-weight:bold;
      line-height:100%;
      margin-top:0;
      margin-right:0;
      margin-bottom:10px;
      margin-left:0;
      text-align:left;
    }

    h4, .h4{
      color:#202020;
      display:block;
      font-family:Arial;
      font-size:22px;
      font-weight:bold;
      line-height:100%;
      margin-top:0;
      margin-right:0;
      margin-bottom:10px;
      margin-left:0;
      text-align:left;
    }

    #templateHeader{
      background-color:#FFFFFF;
      border-bottom:0;
    }

    .headerContent{
      color:#202020;
      font-family:Arial;
      font-size:34px;
      font-weight:bold;
      line-height:100%;
      padding:0;
      text-align:center;
      vertical-align:middle;
    }

    /* Yahoo! Mail Override */
    .headerContent a:link, .headerContent a:visited, .headerContent a .yshortcuts {
      color:#336699;
      font-weight:normal;
      text-decoration:underline;
    }

    #headerImage{
      height:auto;
      max-width:600px !important;
    }

    #templateContainer, .bodyContent{
      background-color:#FFFFFF;
    }

    .bodyContent div{
      color:#505050;
      font-family:Arial;
      font-size:14px;
      line-height:150%;
      text-align:left;
    }

    /* Yahoo! Mail Override */
    .bodyContent div a:link, .bodyContent div a:visited, .bodyContent div a .yshortcuts {
      color:#336699;
      font-weight:normal;
      text-decoration:underline;
    }

    .templateDataTable{
      background-color:#FFFFFF;
      border:1px solid #DDDDDD;
    }

    .dataTableHeading{
      background-color:#222222;
      color:#FFFFFF;
      font-family:Helvetica;
      font-size:14px;
      font-weight:bold;
      line-height:150%;
      text-align:left;
    }

    /* Yahoo! Mail Override */
    .dataTableHeading a:link, .dataTableHeading a:visited, .dataTableHeading a .yshortcuts {
      color:#FFFFFF;
      font-weight:bold;
      text-decoration:underline;
    }

    .dataTableContent{
      border-top:1px solid #DDDDDD;
      border-bottom:0;
      color:#202020;
      font-family:Helvetica;
      font-size:12px;
      font-weight:bold;
      line-height:150%;
      text-align:left;
    }

    /* Yahoo! Mail Override */
    .dataTableContent a:link, .dataTableContent a:visited, .dataTableContent a .yshortcuts {
      color:#202020;
      font-weight:bold;
      text-decoration:underline;
    }

    .templateButton{
      -moz-border-radius:3px;
      -webkit-border-radius:3px;
      background-color:#336699;
      border:0;
      border-collapse:separate !important;
      border-radius:3px;
    }

    /* Yahoo! Mail Override */
    .templateButton, .templateButton a:link, .templateButton a:visited, .templateButton a .yshortcuts {
      color:#FFFFFF;
      font-family:Arial;
      font-size:15px;
      font-weight:bold;
      letter-spacing:-.5px;
      line-height:100%;
      text-align:center;
      text-decoration:none;
    }

    .bodyContent img{
      display:inline;
      height:auto;
    }

    #templateFooter{
      background-color:#FFFFFF;
      border-top:0;
    }

    .footerContent div{
      color:#707070;
      font-family:Arial;
      font-size:12px;
      line-height:125%;
      text-align:center;
    }

    /* Yahoo! Mail Override */
    .footerContent div a:link, .footerContent div a:visited, .footerContent div a .yshortcuts {
      color:#336699;
      font-weight:normal;
      text-decoration:underline;
    }

    .footerContent img{
      display:inline;
    }

    #utility{
      background-color:#FFFFFF;
      border:0;
    }

    #utility div{
      text-align:center;
    }

    #monkeyRewards img{
      max-width:190px;
    }
		</style>
	</head>
  <body leftmargin="0" marginwidth="0" topmargin="0" marginheight="0" offset="0">
    <center>
      <table border="0" cellpadding="0" cellspacing="0" height="100%" width="100%" id="backgroundTable">
        <tr>
          <td align="center" valign="top" style="padding-top:20px;">
            <table border="0" cellpadding="0" cellspacing="0" width="600" id="templateContainer">
              <tr>
                <td align="center" valign="top">
                  <table border="0" cellpadding="0" cellspacing="0" width="600" id="templateHeader">
                    <tr>
                      <td class="headerContent">
                        <img src="https://s3.amazonaws.com/tm.org.mx/tsv/email_header_03.png" style="max-width:600px;" id="headerImage campaign-icon" mc:label="header_image" mc:edit="header_image" mc:allowdesigner mc:allowtext />
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
              <tr>
                <td align="center" valign="top">
                  <table border="0" cellpadding="0" cellspacing="0" width="600" id="templateBody">
                    <tr>
                      <td valign="top">
                        <table border="0" cellpadding="20" cellspacing="0" width="100%">
                          <tr>
                            <td valign="top" class="bodyContent">
                              <div mc:edit="std_content00">
                                <h2 class="h2">{{.Title}}</h2>
                                {{.Content}}
                              </div>
                            </td>
                          </tr>
                          {{if .Items}}
                          <tr>
                            <td valign="top" style="padding-top:0; padding-bottom:0;">
                              <table border="0" cellpadding="10" cellspacing="0" width="100%" class="templateDataTable">
                                <tr>
                                  <th scope="col" valign="top" width="25%" class="dataTableHeading" mc:edit="data_table_heading00">
                                    FECHA
                                  </th>
                                  <th scope="col" valign="top" width="25%" class="dataTableHeading" mc:edit="data_table_heading01">
                                    VALOR
                                  </th>
                                  <th scope="col" valign="top" width="50%" class="dataTableHeading" mc:edit="data_table_heading02">
                                    DESCRIPCION
                                  </th>
                                </tr>
                                {{range .Items}}
                                  <tr mc:repeatable>
                                    <td valign="top" class="dataTableContent" mc:edit="data_table_content00">
                                      {{.Date}}
                                    </td>
                                    <td valign="top" class="dataTableContent" mc:edit="data_table_content01">
                                      {{.Amount}}
                                    </td>
                                    <td valign="top" class="dataTableContent" mc:edit="data_table_content02">
                                      {{.Description}}
                                    </td>
                                  </tr>
                                {{end}}
                              </table>
                            </td>
                          </tr>
                          {{end}}
                          <tr>
                            <td align="center" valign="top" style="padding-top:0;">
                              <br />
                              <table border="0" cellpadding="15" cellspacing="0" class="templateButton">
                                <tr>
                                  <td valign="middle" class="templateButtonContent">
                                    <div mc:edit="std_content02">
                                      <a href="http://www.testigosocial.mx/" target="_blank">Visitar Testigo Social</a>
                                    </div>
                                  </td>
                                </tr>
                              </table>
                            </td>
                          </tr>
                        </table>
                        <!-- // End Module: Standard Content \\ -->

                      </td>
                    </tr>
                  </table>
                  <!-- // End Template Body \\ -->
                </td>
              </tr>
              <tr>
                <td align="center" valign="top">
                  <!-- // Begin Template Footer \\ -->
                  <table border="0" cellpadding="10" cellspacing="0" width="600" id="templateFooter">
                    <tr>
                      <td valign="top" class="footerContent">

                        <!-- // Begin Module: Transactional Footer \\ -->
                        <table border="0" cellpadding="10" cellspacing="0" width="100%">
                          <tr>
                            <td valign="top">
                              <div mc:edit="std_footer">
                                <strong>Transparencia Mexicana AC</strong>
                                <br />
                                Dulce Olivia 73<br />
                                Colonia Villa Coyoacán, Delegación Coyoacán<br />
                                CP 04000<br />
                                México, Ciudad de México
                              </div>
                            </td>
                          </tr>
                          <tr>
                            <td valign="middle" id="utility">
                              <div mc:edit="std_utility">
                                &nbsp;<a href="https://www.twitter.com/testigo_social" target="_blank">Twitter</a> | <a href="https://www.facebook.com/testigosocial">Facebook</a>
                              </div>
                            </td>
                          </tr>
                        </table>
                        <!-- // End Module: Transactional Footer \\ -->

                      </td>
                    </tr>
                  </table>
                  <!-- // End Template Footer \\ -->
                </td>
              </tr>
            </table>
            <br />
          </td>
        </tr>
      </table>
    </center>
  </body>
</html>
`)
