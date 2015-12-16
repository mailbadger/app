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
use newsletters\Exceptions\ModelNotFoundException;
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
     * @param $subject
     * @param null $cc
     * @return mixed
     */
    public function sendEmail(
        SesClient $client,
        $html,
        $email,
        $fromEmail,
        $fromName,
        $subject,
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
                        'Data'    => $html
                    ],
                ],
                'Subject' => [
                    'Charset' => 'UTF-8',
                    'Data'    => $subject,
                ],
            ],
            'Source' => $fromEmail,
        ];

        if (isset($cc)) {
            $data['Destination']['CcAddresses'] = [$cc];
        }

        $response = $client->sendEmail($data);

        return $response->get('MessageId');
    }

    /**
     * Find all sent emails from campaign
     * @param $campaignId
     * @param bool $paginate
     * @param int $perPage
     * @return mixed
     */
    public function findAllSentEmailsByCampaignId($campaignId, $paginate = false, $perPage = 10)
    {
        $emails = $this->sentEmailRepository->with(['bounces', 'complaints'])->scopeQuery(function ($q) use ($campaignId) {
            return $q->where('campaign_id', $id);
        });

        return (!empty($paginate)) ? $emails->paginate($perPage) : $emails->all();
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
     * Increments the sent email opens
     * @param $campaignId
     * @param $subscriberId
     * @return mixed
     */
    public function incrementOpensByToken($token)
    {
        $sentEmail = $this->sentEmailRepository->findByField('token', $token)->first();

        if (empty($sentEmail)) {
            throw new ModelNotFoundException();
        }

        return $this->updateSentEmail(['opens' => $sentEmail->opens + 1], $sentEmail->id);
    }

    /**
     * Generates a random unique token for tracking email opens
     *
     * @param $length
     * @return mixed
     */
    public function generateUniqueToken($length = 32)
    {
        $token = str_random($length);

        $check = $this->sentEmailRepository->findByField('token', $token)->first();

        if (!empty($check)) {
            return $this->generateUniqueToken($length);
        }

        return $token;
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
     * Update sent email
     *
     * @param array $data
     * @param $id
     * @return mixed
     */
    public function updateSentEmail(array $data, $id)
    {
        return $this->sentEmailRepository->update($data, $id);
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
