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
use newsletters\Repositories\BounceRepository;
use newsletters\Repositories\ComplaintRepository;
use Aws\Ses\SesClient;

class EmailService
{
    /**
     * @var SentEmailRepository
     */
    protected $sentEmailRepository;

    /**
     * @var BounceRepository
     */
    protected $bounceRepository;

    /**
     * @var ComplaintRepository
     */
    protected $complaintRepository;

    public function __construct(
        SentEmailRepository $sentEmailRepository, 
        BounceRepository $bounceRepository, 
        ComplaintRepository $complaintRepository
    ) {
        $this->sentEmailRepository = $sentEmailRepository;
        $this->bounceRepository = $bounceRepository;
        $this->complaintRepository = $complaintRepository;
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
     * @return mixed 
     */
    public function sendEmail(
        SesClient $client,
        TemplateService $templateService,
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
            'Destination' => [   
                'ToAddresses' => [$email],
            ],
            'Message' => [ 
                'Body' => [
                    'Html' => [
                        'Charset' => 'UTF-8',
                        'Data' => $templateService->renderTemplate($templateId, $name, $email, $customFields), 
                    ],
                ],
                'Subject' => [ 
                    'Charset' => 'UTF-8',
                    'Data' => $subject, 
                ],
            ],
            'Source' => $fromEmail, 
        ]; 

        if(isset($cc)) {
            $data['Destination']['CcAddresses'] = [$cc];
        }

        $response = $client->sendEmail($data);
 
        return $response->get('MessageId'); 
    }

    /**
     * Find a sent email by the message id
     * @param $messageId
     * @return mixed
     */
    public function findSentEmailByMessageId($messageId)
    {
        return $this->sentEmailRepository->findByField('message_id', $messageId);
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
     * Create bounce
     *
     * @param array $data
     * @return mixed
     */
    public function createBounce(array $data)
    {
        return $this->bounceRepository->create($data);
    }

    /**
     * Create complaint
     *
     * @param array $data
     * @return mixed
     */
    public function createComplaint(array $data)
    {
        return $this->complaintRepository->create($data);
    }
}
