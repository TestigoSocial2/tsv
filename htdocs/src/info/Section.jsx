import React from 'react';
import Description from '../base/Description.jsx';
import Glossary from './Glossary.jsx';

class Section extends React.Component {
  render() {
    return(
      <div>
        <Description
          title="¿Qué es Testigo Social 2.0?"
          color="purple"
          content="Testigo Social 2.0 es una nueva plataforma para que los ciudadanos mejoremos la forma en que el gobierno compra con nuestro dinero, el dinero público. Accede a información sobre los procesos de compra del gobierno y organízate con otros ciudadanos para proponer ajustes en los proyectos que te importan o revisar que los recursos asignados se gasten de forma correcta. En esta sección encontrarás distintos recursos e información para hacer un mejor uso de las herramientas disponibles en Testigo Social 2.0." />
        <Glossary />
      </div>
    );
  }
}

export default Section;
