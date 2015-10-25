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
use Mail;

class EmailService
{
    /**
     * @var SentEmailRepository
     */
    protected $sentEmailRepository;

    public function __construct(SentEmailRepository $sentEmailRepository)
    {
        $this->sentEmailRepository = $sentEmailRepository;
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
        'email'         => $email,
        'template_id'   => $templateId,
        'custom_fields' => $customFields,
    ];

    Mail::send('emails.template', $data,
        function ($message) use ($email, $fromEmail, $fromName, $subject, $cc) {
            $message->from($fromEmail, $fromName);
            $message->to($email)->subject($subject);

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

    /**
     * Sets the SES user configuration
     *
     * @param $key
     * @param $secret
     * @param $region
     */
    public function setSesConfig($key, $secret, $region)
    {
        if (empty($key) || empty($secret) || empty($region)) {
            throw new InvalidArgumentException('SES configuration is not set.');
        }

        config([
            'services.ses.key'    => $key,
            'services.ses.secret' => $secret,
            'services.ses.region' => $region,
        ]);
    }

}
