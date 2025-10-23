package application

import (
	"context"
	"fmt"

	"github.com/Maxito7/hotel_backend/internal/domain"
	"github.com/Maxito7/hotel_backend/internal/email"
	"github.com/Maxito7/hotel_backend/internal/infrastructure/repository"
)

type ContactService struct {
	repo        repository.ContactRepository
	emailClient *email.Client
}

func NewContactService(r repository.ContactRepository, emailClient *email.Client) *ContactService {
	return &ContactService{
		repo:        r,
		emailClient: emailClient,
	}
}

func (s *ContactService) Create(ctx context.Context, req domain.CreateContactRequest) (int64, error) {
	// Crear el contacto en la base de datos
	id, err := s.repo.Create(ctx, req)
	if err != nil {
		return 0, err
	}

	// Enviar email de notificación si el cliente está configurado
	if s.emailClient != nil {
		if err := s.enviarEmailContacto(req); err != nil {
			// Log error pero no fallar
			fmt.Printf("Error al enviar email de contacto: %v\n", err)
		}
	}

	return id, nil
}

// enviarEmailContacto envía un email cuando alguien completa el formulario de contacto
func (s *ContactService) enviarEmailContacto(req domain.CreateContactRequest) error {
	subject := fmt.Sprintf("Nuevo mensaje de contacto - %s", req.Nombre)

	// Manejar valores opcionales (punteros)
	telefono := "No proporcionado"
	if req.Telefono != nil && *req.Telefono != "" {
		telefono = *req.Telefono
	}

	mensaje := "Sin mensaje"
	if req.Mensaje != nil && *req.Mensaje != "" {
		mensaje = *req.Mensaje
	}

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="es">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Nuevo Mensaje de Contacto</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4;">
	<table width="100%%" cellpadding="0" cellspacing="0" style="background-color: #f4f4f4; padding: 20px;">
		<tr>
			<td align="center">
				<table width="600" cellpadding="0" cellspacing="0" style="background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
					<!-- Header -->
					<tr>
						<td style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 30px 20px; text-align: center;">
							<h1 style="color: #ffffff; margin: 0; font-size: 24px;">📬 Nuevo Mensaje de Contacto</h1>
						</td>
					</tr>

					<!-- Contenido -->
					<tr>
						<td style="padding: 40px 30px;">
							<div style="background-color: #f8f9fa; border-left: 4px solid #667eea; padding: 20px; margin-bottom: 30px;">
								<h2 style="margin: 0 0 15px 0; color: #333; font-size: 18px;">Información del Cliente</h2>
								<table width="100%%" cellpadding="0" cellspacing="0">
									<tr>
										<td style="padding: 8px 0;"><strong>Nombre:</strong></td>
										<td style="padding: 8px 0; text-align: right;">%s</td>
									</tr>
									<tr>
										<td style="padding: 8px 0;"><strong>Email:</strong></td>
										<td style="padding: 8px 0; text-align: right;">%s</td>
									</tr>
									<tr>
										<td style="padding: 8px 0;"><strong>Teléfono:</strong></td>
										<td style="padding: 8px 0; text-align: right;">%s</td>
									</tr>
								</table>
							</div>

							<!-- Mensaje -->
							<h3 style="color: #333; margin-bottom: 15px;">Mensaje:</h3>
							<div style="background-color: #f8f9fa; padding: 20px; border-radius: 8px; border: 1px solid #e0e0e0;">
								<p style="margin: 0; color: #555; line-height: 1.6;">%s</p>
							</div>

							<!-- Información adicional -->
							<div style="margin-top: 30px; padding: 20px; background-color: #fff3cd; border-radius: 8px; border-left: 4px solid #ffc107;">
								<p style="margin: 0; color: #856404; font-size: 14px;">
									💡 <strong>Responder pronto:</strong> Este cliente está esperando una respuesta. 
									Contacta directamente a <a href="mailto:%s" style="color: #667eea;">%s</a>
								</p>
							</div>
						</td>
					</tr>

					<!-- Footer -->
					<tr>
						<td style="background-color: #f8f9fa; padding: 20px; text-align: center; border-top: 1px solid #e0e0e0;">
							<p style="margin: 0; color: #999; font-size: 12px;">
								Este es un correo automático del sistema de contacto de Hotel Inca
							</p>
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</table>
</body>
</html>
	`,
		req.Nombre,
		req.Email,
		telefono,
		mensaje,
		req.Email,
		req.Email,
	)

	// Enviar a la empresa (hotelinca.reservas@gmail.com está configurado en SMTP_FROM_EMAIL)
	return s.emailClient.SendEmail("hotelinca.reservas@gmail.com", subject, htmlBody)
}

func (s *ContactService) List(ctx context.Context) ([]domain.Contact, error) {
	return s.repo.List(ctx)
}

func (s *ContactService) UpdateEstado(ctx context.Context, id int64, estado domain.EstadoFormulario) error {
	return s.repo.UpdateEstado(ctx, id, estado)
}
