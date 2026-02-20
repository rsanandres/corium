import { NextRequest, NextResponse } from 'next/server';

// Contact form is disabled in this environment (no email/CAPTCHA credentials).
// To re-enable, uncomment the original handler below and configure
// RECAPTCHA_SECRET_KEY, EMAIL_USER, and EMAIL_PASS environment variables.

export async function POST(_req: NextRequest) {
  return NextResponse.json(
    { error: 'Contact form is not available in this environment.' },
    { status: 503 }
  );
}

/*
import nodemailer from 'nodemailer';

async function verifyCaptcha(token: string) {
  try {
    const response = await fetch('https://www.google.com/recaptcha/api/siteverify', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: `secret=${process.env.RECAPTCHA_SECRET_KEY}&response=${token}`,
    });

    const data = await response.json();
    return data.success;
  } catch (error) {
    console.error('CAPTCHA verification error:', error);
    return false;
  }
}

export async function POST(req: NextRequest) {
  try {
    const { name, email, message, captchaToken } = await req.json();

    if (!name || !email || !message) {
      return NextResponse.json({ error: 'All fields are required.' }, { status: 400 });
    }

    if (!captchaToken) {
      return NextResponse.json({ error: 'CAPTCHA verification failed.' }, { status: 400 });
    }

    const isValidCaptcha = await verifyCaptcha(captchaToken);
    if (!isValidCaptcha) {
      return NextResponse.json({ error: 'CAPTCHA verification failed.' }, { status: 400 });
    }

    const transporter = nodemailer.createTransport({
      service: 'gmail',
      auth: {
        user: process.env.EMAIL_USER,
        pass: process.env.EMAIL_PASS,
      },
    });

    const emailHtml = `
      <div style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
        <h2 style="color: #4A90E2;">Contact Form Submission</h2>
        <table style="width: 100%; border-collapse: collapse; margin-top: 20px;">
          <tr style="background-color: #f2f2f2;">
            <td style="padding: 12px; border: 1px solid #ddd; font-weight: bold;">Field</td>
            <td style="padding: 12px; border: 1px solid #ddd;">Submission</td>
          </tr>
          <tr>
            <td style="padding: 12px; border: 1px solid #ddd; font-weight: bold;">Name:</td>
            <td style="padding: 12px; border: 1px solid #ddd;">${name}</td>
          </tr>
          <tr>
            <td style="padding: 12px; border: 1px solid #ddd; font-weight: bold;">Email:</td>
            <td style="padding: 12px; border: 1px solid #ddd;">${email}</td>
          </tr>
        </table>
        <div style="margin-top: 20px; padding: 15px; background-color: #f9f9f9; border-left: 4px solid #4A90E2;">
          <h3 style="margin-top: 0; color: #4A90E2;">Message:</h3>
          <p style="margin-bottom: 0;"><em>"${message}"</em></p>
        </div>
        <hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
        <p style="font-size: 0.8em; color: #888;">Somebody contacted you! Exciting!</p>
      </div>
    `;

    const mailOptions = {
      from: `Resume Contact <${process.env.EMAIL_USER}>`,
      to: 'contact@rsanandres.com',
      subject: `New Resume Inquiry from: ${name}`,
      text: `New submission from ${name} (${email}): ${message}`,
      html: emailHtml,
    };

    try {
      await transporter.sendMail(mailOptions);
      return NextResponse.json({
        success: true,
        message: 'Form submitted successfully!'
      });
    } catch (emailError) {
      console.error('Email sending error:', emailError);
      return NextResponse.json({ error: 'Failed to send email.' }, { status: 500 });
    }
  } catch (error) {
    console.error('Contact form error:', error);
    return NextResponse.json({ error: 'Failed to process form submission.' }, { status: 500 });
  }
}
*/
