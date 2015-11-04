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

        if(isset($cc)) {
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
            return $q->whereHas('campaign', function ($q) use ($campaignId) {
                return $q->where('id', $campaignId);
            });
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
     * Find campaign reports for sent emails
     *
     * @param $campaignId
     * @return mixed
     */
    public function findSendsReportByCampaignId($campaignId)
    {
        $complaints = $bounces = $sent = $opens = 0;

        //TODO fetch sent emails by chunks
        $this->sentEmailRepository->with(['complaintsCount', 'bouncesCount'])
            ->findByField('campaign_id', $campaignId)
            ->each(function ($e) use(&$complaints, &$bounces, &$sent, &$opens) { 
                $sent++;
                $opens += $e->opens;
                $complaints += $e->complaintsCount->sum('complaints'); 
                $bounces += $e->bouncesCount->sum('bounces');
            });

        return ['bounces' => $bounces, 'complaints' => $complaints, 'sent' => $sent, 'opens' => $opens];  
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
