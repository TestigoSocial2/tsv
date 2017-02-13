import React from 'react';

class Glossary extends React.Component {
  componentDidMount() {
    $('#glossaryContents td div.content').hide();
    $('#glossaryContents td').on( 'click', function(e){
      $(e.currentTarget).find('div.content').slideToggle();
    });
  }

  render() {
    return(
      <div className="row">
        <div className="inner-row">
          <h1>Glosario</h1>
          <table id="glossaryContents" className="table table-striped green">
            <tbody>
              <tr>
                <td>
                  <h4>Adjudicación Directa</h4>
                  <div className="content">
                    <p>Proceso de contratación que por la naturaleza del bien a comprar, o servicio a contratar, u obra pública a realizar por el que la entidad convocante celebra un contrato con una persona en específico siempre y cuando no exceda los montos máximos establecidos por el Presupuesto de Egresos de la Federación y cumplea con las condiciones establecidas en la ley.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Auditoría Superior de la Federación</h4>
                  <div className="content">
                    <p>Ente que depende del poder legsilativo que administra y fiscaliza los recursos públicos federales.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Bases de la Licitación</h4>
                  <div className="content">
                    <p>Son los lineamientos o requisitos específicos que establecen el objeto de la licitación y las condiciones de participación.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>CompraNet</h4>
                  <div className="content">
                    <p>Sistema electrónico de información pública gubernamental sobre procedimientos de contratación pública. Es de consulta gratuita y es el medio por el cual se desarrollan procedimientos de contratación.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Compras Públicas</h4>
                  <div className="content">
                    <p>Adquisición y arrendamiento de todo tipo de bienes, contratación de servicios y construcción de obra que se realicen con recursos del Estado.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Contratación Pública</h4>
                  <div className="content">
                    <p>Es el vínculo jurídico entre un ente público y una contraparte que puede ser pública o privada para adquirir o arrendar cualquier tipo de bienes, contratar servicios o realizar obra pública y abarca desde la firma del contrato hasta la evaluación de cumplimiento del contrato.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Contratista</h4>
                  <div className="content">
                    <p>Persona que celebra un contrato público.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Contrato</h4>
                  <div className="content">
                    <p>Acuerdo formal o vinculo juridico en el que establecen los derechos y obligaciones de las partes entre una entidad pública (convocante) y un particular para la adquisición o arrendamiento todo tipo de bienes, contratación de servicios,  o realización de obra pública. (en el contexto de compras públicas).</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Convenio modificatorio</h4>
                  <div className="content">
                    <p>Acuerdo que busca modificar una parte del contrato en el cual se podrán modificar las condiciones del contrato original, por ejemplo el monto y los plazos, siguiendo con una serie de condiciones legales.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Datos Abiertos</h4>
                  <div className="content">
                    <p>Datos abiertos son datos digitales que son puestos a disposición con las características técnicas y jurídicas necesarias para que puedan ser usados, reutilizados y redistribuidos libremente por cualquier persona, en cualquier momento y en cualquier lugar.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Estándar de Datos de Contrataciones Abiertas</h4>
                  <div className="content">
                    <p>Es un estándar de datos abiertos para la publicación de información estructurada en todas las etapas de un proceso de contratación: desde la planeación hasta la implementación.La publicación de datos del OCDS puede posibilitar mayor transparencia en las contrataciones públicas y puede facilitar un análisis accesible y profundo de la eficiencia, efectividad, legitimidad e integridad de los sistemas de contrataciones públicas.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Etapas de la Contratación</h4>
                  <div className="content">
                    <p>Existen varias etapas en un proceso de contratación:</p>
                    <p><strong>Planeación:</strong> En relación a la contratación, se refiere a una etapa preliminar de preparación (Presupuestos, planes de proyecto, planes de contratación, estudios de mercado, información de audiencias públicas).</p>
                    <p><strong>Licitación:</strong> Proceso de convocatoria para que se escojan las mejores condiciones de compra y se presenten proposiciones para procedimientos de contratación, adquisición, arrendamiento, o prestación de servicios. (Anuncios de licitación,especificaciones,partidas,importes,consultas).</p>
                    <p><strong>Adjudicación:</strong> Atribuir o declarar un contrato al ganador del concurso de licitación (Detalles de la adjudicación, información sobre el licitante, evaluación de la licitación, importes).</p>
                    <p><strong>Contratación:</strong> Proceso por el cual se compromete a un individuo y organización por medio de contrato a desempeñar las obligaciones estipuladas en el mismo. (Detalles finales, contrato firmado, enmiendas, importes).</p>
                    <p><strong>Implementación:</strong> Etapa final de la contratación, se refiere a la ejecución del contrato. (Pagos, actualizaciones de progreso, ubicación, extensiones, enmiendas, información sobre la Realización o Terminación).</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Evaluación de propuestas</h4>
                  <div className="content">
                    <p>Análisis de las propuestas recibidas con base en los requisitos y criterios establecidos en la convocatoria o en las bases de la licitación para determinar su solvencia técnica y económica.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Fallo</h4>
                  <div className="content">
                    <p>Documento emitido por el convocante en el que se anuncia el nombre del licitante al que se adjudicará el contrato; las propuestas desechadas con las razones por las cuales lo fueron; fecha, lugar y hora para la firma del contrato, garantías; nombre cargo y firma del servidor público que lo emite.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Invitación a cuando menos tres</h4>
                  <div className="content">
                    <p>Proceso de contratación que por la naturaleza del bien a comprar, o servicio a contratar, y obra pública a realizar difunde la invitación en CompraNet y en la página de internet de la dependencia o entidad a cuando menos tres personas que tengan las capacidades establacidas en la ley.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Junta de aclaraciones</h4>
                  <div className="content">
                    <p>Es una reunión convocada por la entidad contratante en la que los concursantes pueden hacer preguntas sobre la convocatoria. De acuerdo a la legislación mexicana, se deberá realizar al menos una junta de aclaraciones a la convocatoria de la licitación.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Ley de adquisiciones, arrendamientos y servicios del sector público</h4>
                  <div className="content">
                    <p>Tiene como objetivo reglamentar la aplicación del artículo 134 de la Constitución Política de los Estados Unidos Mexicanos. El cual se refiere a que la administración de los recursos disponibles se realice con eficiencia, eficacia, economía, transparencia y honradez para satisfacer los objetivos a los que estén destinados, en materia de las adquisiciones, arrendamientos de bienes inmuebles y prestación de servicios de cualquier naturaleza.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Ley de obras públicas y servicios relacionados con las mismas</h4>
                  <div className="content">
                    <p>Tiene como objetivo reglamentar la aplicación del artículo 134 de la Constitución Política de los Estados Unidos Mexicanos. El cual se refiere a que la administración de los recursos disponibles se realice con eficiencia, eficacia, economía, transparencia y honradez para satisfacer los objetivos a los que estén destinados, en materia de contrataciones, obras públicas y servicios relacionados con las mismas.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Licitación pública</h4>
                  <div className="content">
                    <p>Procedimiento de contratación mediante convocatoria pública para que libremente se presenten proposiciones para la contratación de obra, arrendamiento, adquisiciones, y servicios.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Licitante</h4>
                  <div className="content">
                    <p>Persona que participa en cualquier procedimiento de contratación pública o invitación a cuando menos tres personas.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Método de evaluación</h4>
                  <div className="content">
                    <p>Criterio que será utilizado por la entidad convocante para determinar la solvencia de las propuestas presentadas por los interesados, asi como para determinar a la propuesta que será adjudicada. El metodo de evaluación puede ser binario, por puntos y porcentajes o por costo-beneficio.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Pre-bases de la licitación</h4>
                  <div className="content">
                    <p>Es el borrador o versión preliminar de las bases y se difunden para ser sometidas a la opinión de los interesados, quienes tienen la oportunidad de formular comentarios al documento. Estas no son vinculantes.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Presentación de propuestas</h4>
                  <div className="content">
                    <p>Etapa en la que los interesados en el procedimiento de contratación presentan propuestas para cumplir con los requisitos técnicos y objetivos con el objeto de resultar adjudicados.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Proceso de Contratación</h4>
                  <div className="content">
                    <p>Es la secuencia de eventos orientada a la adjudicaciín de contrato. Iniciando con la planeación y termina con la firma.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Secretaría de la Función Pública</h4>
                  <div className="content">
                    <p>Entidad pública que determina la política de compras de la Federación (compras públicas) lo que significa que está facultada para interpretar la Ley de Obra Pública y la Ley de Adquisiciones, Arrendamientos y Servicios. Además, audita el gasto de los recursos federales, y vigila el desempeño de los servidores públicos federales, entre otras funciones.</p>
                  </div>
                </td>
              </tr>
              <tr>
                <td>
                  <h4>Testigo Social</h4>
                  <div className="content">
                    <p>Es una figura  creada por la Sociedad Civil que tiene el objetivo de observar de manera independiente un proceso de contratación. Esta figura ha sido recogida en diversas leyes del Sistema Jurídico Mexicano. Los testigos sociales son designados por la Secretaría de Función Pública para que acompañan y participan en todos los procesos de licitación pública cuyos montos excedan cierta cantidad. Los testigos sociales participan en todas las etapas del proceso de contratación y emiten un testimonio final que incluye observaciones y/o recomendaciones. En caso de detectar una irregularidad en el proceso, lo remite a las autorridades correspondientes.</p>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    );
  }
}

export default Glossary;
