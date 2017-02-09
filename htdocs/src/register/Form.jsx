import React from 'react';

class Form extends React.Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    let agencies = $('div.agency-grid div.item');
    let selectedAgencies = $('div.agency-grid input#selectedAgencies');
    agencies.click(function() {
      // Toggle state
      let a = $(this);
      let lbl = a.find('span.label');
      a.toggleClass('selected');
      if( lbl.text().toLowerCase() == 'seguir' ) {
        lbl.text( 'siguiendo' );
      } else {
        lbl.text( 'seguir' );
      }

      // Adjust result
      let res = [];
      agencies.each( function( i, v ) {
        if( $(v).hasClass( 'selected') ) {
          res.push( $(v).data('value') );
        }
      });
      selectedAgencies.val(JSON.stringify(res));
    });

    let projects = $('div.project-grid div.item');
    let selectedProjects = $('div.project-grid input#selectedProjects');
    projects.click(function() {
      // Toggle state
      let p = $(this);
      let lbl = p.find('span.label');
      p.toggleClass('selected');
      if( lbl.text().toLowerCase() == 'seguir' ) {
        lbl.text( 'siguiendo' );
      } else {
        lbl.text( 'seguir' );
      }

      // Adjust result
      let res = [];
      projects.each( function( i, v ) {
        if( $(v).hasClass( 'selected') ) {
          res.push( $(v).data('value') );
        }
      });
      selectedProjects.val(JSON.stringify(res));
    });

    let form = $('form#registerForm');
    form.pixativeFormValidator({
      msgHolder: 'title',
      errorMessages: {
        required: 'El campo es requerido',
        minlength: 'El valor proporcionado debe ser de al menos %s caracteres',
        maxlength: 'El valor proporcionado debe ser de máximo %s caracteres',
        email: 'El valor porporcionado no parece ser una dirección de correo valida',
        integer: 'El valor proporcionado debe ser un número',
        phone: 'El valor porporcionado no parece ser un número telefonico valido'
      },
      onError: function() {
        $('.validator-error').tooltip('destroy');
        $('.validator-error').tooltip();
      },
      onSuccess: () => {
        // Prepare data
        let data = {}
        form.serializeArray().forEach(function(el) {
          if( el.value == 'on' ) {
            el.value = true;
          }
          if( el.value == 'off' ) {
            el.value = false;
          }
          data[el.name] = el.value;
        });
        data.selectedAgencies = JSON.parse(data.selectedAgencies);
        data.selectedProjects = JSON.parse(data.selectedProjects);
        this.props.onSubmit(data);
      }
    });
  }

  render() {
    return(
      <div id="setup">
        <form id="registerForm">
          <div className="row">
            <div className="col-md-3">
              <h2 className="number-title"><span>1</span>Registro</h2>
              <p className="txt-italic">Completa la siguiente información</p>
            </div>
            <div className="col-md-9">
              <table className="table table-striped green">
                <tbody>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Correo Electrónico</label></td>
                    <td className="col-md-4">
                      <input type="email"
                             className="form-control"
                             id="user"
                             name="user"
                             placeholder="Email"
                             data-validator-required="true"
                             data-validator-email="true" />
                    </td>
                    <td className="txt-italic">La cuenta de correo electrónico registrada será tu nombre de usuario</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Contraseña</label></td>
                    <td className="col-md-4">
                      <input type="password"
                             className="form-control"
                             id="password"
                             name="password"
                             placeholder="Contraseña"
                             data-validator-required="true"
                             data-validator-minlength="6" />
                    </td>
                    <td className="txt-italic">La contraseña deberá tener una extensión de al menos 6 caracteres y contener al menos un número</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Tipo de Usuario</label></td>
                    <td className="col-md-4">
                      <select id="userType" name="userType" className="form-control" data-validator-required="true">
                        <option value="-" disabled="disaled">Selecciona una de las siguientes opciones</option>
                        <option value="public">Ciudadano</option>
                        <option value="legal">Legislador</option>
                        <option value="gov">Funcionario Público</option>
                        <option value="media">Periodista</option>
                        <option value="startup">Emprendedor</option>
                        <option value="business">Empresario</option>
                        <option value="other">Otro</option>
                      </select>
                    </td>
                    <td className="txt-italic">Selecciona una de las siguientes opciones</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Edad</label></td>
                    <td className="col-md-4">
                      <select id="age" name="age" className="form-control" data-validator-required="true">
                        <option value="15">15</option>
                        <option value="16">16</option>
                        <option value="17">17</option>
                        <option value="18">18</option>
                        <option value="19">19</option>
                        <option value="20">20</option>
                        <option value="21">21</option>
                        <option value="22">22</option>
                        <option value="23">23</option>
                        <option value="24">24</option>
                        <option value="25">25</option>
                        <option value="26">26</option>
                        <option value="27">27</option>
                        <option value="28">28</option>
                        <option value="29">29</option>
                        <option value="30">30</option>
                        <option value="31">31</option>
                        <option value="32">32</option>
                        <option value="33">33</option>
                        <option value="34">34</option>
                        <option value="35">35</option>
                        <option value="36">36</option>
                        <option value="37">37</option>
                        <option value="38">38</option>
                        <option value="39">39</option>
                        <option value="40">40</option>
                        <option value="41">41</option>
                        <option value="42">42</option>
                        <option value="43">43</option>
                        <option value="44">44</option>
                        <option value="45">45</option>
                        <option value="46">46</option>
                        <option value="47">47</option>
                        <option value="48">48</option>
                        <option value="49">49</option>
                        <option value="50">50</option>
                        <option value="51">51</option>
                        <option value="52">52</option>
                        <option value="53">53</option>
                        <option value="54">54</option>
                        <option value="55">55</option>
                        <option value="56">56</option>
                        <option value="57">57</option>
                        <option value="58">58</option>
                        <option value="59">59</option>
                        <option value="60">60</option>
                        <option value="61">61</option>
                        <option value="62">62</option>
                        <option value="63">63</option>
                        <option value="64">64</option>
                        <option value="65">65</option>
                        <option value="66">66</option>
                        <option value="67">67</option>
                        <option value="68">68</option>
                        <option value="69">69</option>
                        <option value="70">70</option>
                        <option value="71">71</option>
                        <option value="72">72</option>
                        <option value="73">73</option>
                        <option value="74">74</option>
                        <option value="75">75</option>
                        <option value="76">76</option>
                        <option value="77">77</option>
                        <option value="78">78</option>
                        <option value="79">79</option>
                        <option value="80">80</option>
                        <option value="81">81</option>
                        <option value="82">82</option>
                        <option value="83">83</option>
                        <option value="84">84</option>
                        <option value="85">85</option>
                        <option value="86">86</option>
                        <option value="87">87</option>
                        <option value="88">88</option>
                        <option value="89">89</option>
                        <option value="90">90</option>
                        <option value="91">91</option>
                        <option value="92">92</option>
                        <option value="93">93</option>
                        <option value="94">94</option>
                        <option value="95">95</option>
                        <option value="96">96</option>
                        <option value="97">97</option>
                        <option value="98">98</option>
                        <option value="99">99</option>
                      </select>
                    </td>
                    <td className="txt-italic">Selecciona tu edad (solo aplica para periodista, emprendedor, legislador, funcionario público, ciudadano)</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Código Postal</label></td>
                    <td className="col-md-4">
                      <input type="text"
                             className="form-control"
                             id="postalCode"
                             name="postalCode"
                             data-validator-required="true"
                             data-validator-integer="true" />
                    </td>
                    <td className="txt-italic">Introduce tu código postal</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          <div className="row">
            <div className="col-md-3">
              <h2 className="number-title"><span>2</span>Proyectos</h2>
              <p className="txt-italic">Selecciona las agencias, proyectos o sectores que te interesan</p>
            </div>
            <div className="row">
              <div className="agency-grid col-md-4">
                <h4 className="txt-upper">Agencias</h4>
                <input type="text" id="selectedAgencies" name="selectedAgencies" value="[]" className="hidden" />
                <div className="item" data-value="cdmx">
                  <div className="logo">
                    <img src="https://twitter.com/gobcdmx/profile_image?size=original" />
                  </div>
                  <span className="label">Seguir</span>
                </div>
                <div className="item" data-value="a2">
                  <div className="logo">
                    <img src="https://twitter.com/nvoaeropuertomx/profile_image?size=original" />
                  </div>
                  <span className="label">Seguir</span>
                </div>
                <div className="clear"></div>
              </div>
              <div className="project-grid col-md-3">
                <h4 className="txt-upper">Proyectos</h4>
                <input type="text" id="selectedProjects" name="selectedProjects" value="[]" className="hidden" />
                <div className="item" data-value="p1">
                  <h5>NAICM</h5>
                  <div className="img">
                    <img src="https://twitter.com/nvoaeropuertomx/profile_image?size=original" />
                  </div>
                  <span className="label">Seguir</span>
                  <p>Construcción del nuevo aeropuerto internacional de la ciduad de México.</p>
                </div>
              </div>
              <div className="col-md-2">
                <h4 className="txt-upper">Sectores</h4>
                <ul>
                  <li>
                    <span className="glyphicon glyphicon-home" aria-hidden="true"></span>Infraestructura
                  </li>
                </ul>
              </div>
            </div>
          </div>
          <div className="row">
            <div className="col-md-3">
              <h2 className="number-title"><span>3</span>Contacto</h2>
              <p className="txt-italic">Selecciona los canales por los que Testigo Social 2.0 puede enviarte notificaciones e información.</p>
            </div>
            <div className="col-md-9">
              <table className="table table-striped blue">
                <tbody>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Correo Electrónico</label></td>
                    <td className="col-md-4">
                      <input type="email"
                             className="form-control"
                             id="notificationEmail"
                             name="notificationEmail"
                             data-validator-email="true" />
                    </td>
                    <td>
                      <input id="enableEmailNotifications"
                             name="enableEmailNotifications"
                             type="checkbox"
                             checked="checked" />
                    </td>
                    <td className="txt-italic">Quiero recibir notificaciones a través de un correo electrónico</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>SMS</label></td>
                    <td className="col-md-4">
                      <input type="text"
                             className="form-control"
                             id="notificationSMS"
                             name="notificationSMS"
                             data-validator-integer="true"
                             data-validator-minlength="10"
                             data-validator-maxlength="10" />
                    </td>
                    <td>
                      <input id="enableSMSNotifications"
                             name="enableSMSNotifications"
                             type="checkbox"
                             checked="checked" />
                    </td>
                    <td className="txt-italic">Quiero recibir notificaciones a través de mensajes a un número celular (sin costo)</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Facebook</label></td>
                    <td className="col-md-4">
                      <input disabled="disabled" type="text" className="form-control" id="notificationFB" />
                    </td>
                    <td>
                      <input disabled="disabled" id="enableFB" type="checkbox" />
                    </td>
                    <td className="txt-italic">Quiero recibir notificaciones e información a través mensajes en Facebook Messenger</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Twitter</label></td>
                    <td className="col-md-4">
                      <input disabled="disabled" type="text" className="form-control" id="notificationTW" />
                    </td>
                    <td>
                      <input disabled="disabled" id="enableTW" type="checkbox" />
                    </td>
                    <td className="txt-italic">Quiero recibir notificaciones a través de mensajes directos y notificaciones en Twitter</td>
                  </tr>
                </tbody>
              </table>
              <p>
                Consulta nuestro aviso de privacidad <a href="#" data-toggle="modal" data-target="#privacyNotice">aquí.</a>
              </p>
              <div className="checkbox">
                <label>
                  <input id="acceptPrivacyTerms" type="checkbox" checked="checked" /> Acepto las politicas de privacidad
                </label>
              </div>
              <button type="submit" className="btn btn-black btn-lg">Comienza a recibir notificaciones</button>
            </div>
          </div>
        </form>
      </div>
    );
  }
}

export default Form;
