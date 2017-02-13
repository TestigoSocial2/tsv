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
          content="Testigo Social 2.0 es una nueva plataforma para que los ciudadanos transformemos la forma en que el gobierno compra con nuestro dinero, el dinero público. Infórmate sobre cómo está comprando el gobierno y organízate con otros ciudadanos para cambiar cómo compra el gobierno en los proyectos que te importan." />

        <Glossary />
      </div>
    );
  }
}

export default Section;
