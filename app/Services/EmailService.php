<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 15.9.15
 * Time: 21:06
 */

namespace newsletters\Services;

use Illuminate\Contracts\Mail\Mailer;
use newsletters\Repositories\SentEmailRepository;

class EmailService
{
    /**
     * @var SentEmailRepository
     */
    protected $sentEmailRepository;

    /**
     * @var Mailer
     */
    protected $mailer;

    public function __construct(SentEmailRepository $sentEmailRepository, Mailer $mailer)
    {
        $this->sentEmailRepository = $sentEmailRepository;
        $this->mailer = $mailer;
    }

    /**
     * Send email
     *
     * @param $email
     * @param $name
     * @param $fromEmail
     * @param $fromName
     * @param $subject
     * @param $templateId
     * @param array $customFields
     * @param null $cc
     */
    public function sendEmail(
        $email,
        $name,
        $fromEmail,
        $fromName,
        $subject,
        $templateId,
        $customFields = [],
        $cc = null
    ) {
        $data = [
            'name'          => $name,
            'template_id'   => $templateId,
            'custom_fields' => $customFields,
        ];

        $this->mailer->queueOn('emails', 'emails.template', $data,
            function ($message) use ($email, $fromEmail, $fromName, $subject, $cc) {
                $message->from($fromEmail, $fromName);
                $message->to($email);

                if (isset($cc)) {
                    $message->cc($cc);
                }
            });
    }

    /**
     * Create sent email
     *
     * @param array $data
     * @return mixed
     */
    public function createSentEmail(array $data)
    {
        return $this->sentEmailRepository->create($data);
    }
}