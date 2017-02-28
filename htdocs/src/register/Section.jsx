import React from 'react';
import Description from '../base/Description.jsx';
import Form from './Form.jsx';

class Section extends React.Component {
  constructor(props) {
    super(props);
  }

  processRegister(data) {
    $.ajax({
      type: "POST",
      url: '/profile',
      data: {
        profile: JSON.stringify(data)
      },
      success: function( res ) {
        alert('Tu usuario ha quedado registrado exitosamente en Testigo Social Virtual 2.0');
      }
    });
  }

  render() {
    return (
      <div className="register-description">
        <Description
          title="Notificaciones"
          color="yellow"
          content="Testigo Social 2.0 te puede hacer llegar datos e información específica sobre procedimientos de contratación pública que están en marcha. Desde un aviso de inicio de un nuevo procedimiento hasta la liga para consultar un contrato. Completa la información correspondiente y comienza a recibir notificaciones." />
        <Form onSubmit={this.processRegister} />
      </div>
    );
  }
}

export default Section;
