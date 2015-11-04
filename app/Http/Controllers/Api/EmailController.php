<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;
use newsletters\Http\Requests;
use newsletters\Http\Controllers\Controller;
use Illuminate\Support\Facades\Log;
use Aws\Sns\Message;
use Aws\Sns\MessageValidator;
use GuzzleHttp\Client;
use newsletters\Services\EmailService;
use newsletters\Services\SubscriberService;

class EmailController extends Controller
{
    /**
     * @var EmailService 
     */
    private $service;

    public function __construct(EmailService $service)
    { 
        $this->service = $service;
    }

    public function bounces(MessageValidator $validator, SubscriberService $subscriberService)
    {
        try {
            $message = Message::fromRawPostData();
            $validator->validate($message);
        } catch(Exception $e) {
            abort(404);
        }

        if ('SubscriptionConfirmation' === $message['Type']) {
            (new Client)->get($message['SubscribeURL']);
        } else {
            $bounce = json_decode($message['Message'], true);

            $sentEmail = $this->service->findSentEmailByMessageId($bounce['mail']['messageId'])->first();
            
            if(empty($sentEmail)) {
                abort(404); //If a sent email doesn't exist with that message id, don't write the bounce
            }

            foreach($bounce['bounce']['bouncedRecipients'] as $recipient) {
                $this->service->createBounce([
                    'recipient'     => $recipient['emailAddress'],
                    'sender'        => $bounce['mail']['source'],
                    'action'        => $recipient['action'],
                    'type'          => $bounce['bounce']['bounceType'],
                    'sub_type'      => $bounce['bounce']['bounceSubType'],
                    'timestamp'     => $bounce['bounce']['timestamp'],
                    'sent_email_id' => $sentEmail->id,
                ]);
                
                $subscriber = $subscriberService->findSubscriberByEmail($recipient['emailAddress'])->first();

                if (!empty($subscriber)) {
                    $subscriberService->updateSubscriber(['blacklisted' => true], $subscriber->id);
                }
            }
        }  
    }

    public function complaints(MessageValidator $validator)
    {
        try {
            $message = Message::fromRawPostData();
            $validator->validate($message);
        } catch(Exception $e) {
            abort(404);
        }

        if ('SubscriptionConfirmation' === $message['Type']) {
            (new Client)->get($message['SubscribeURL']);
        } else {
            $complaint = json_decode($message['Message'], true);

            $sentEmail = $this->service->findSentEmailByMessageId($complaint['mail']['messageId'])->first();
 
            if(empty($sentEmail)) {
                abort(404); //If a sent email doesn't exist with that message id, don't write the bounce
            }

            foreach($complaint['complaint']['complainedRecipients'] as $recipient) {
                $this->service->createComplaint([
                    'recipient' => $recipient['emailAddress'], 
                    'sender'    => $complaint['mail']['source'],
                    'type'      => $complaint['complaint']['complaintFeedbackType'], 
                    'timestamp' => $complaint['complaint']['timestamp'],
                    'sent_email_id' => $sentEmail->id,
                ]);
                
                $subscriber = $subscriberService->findSubscriberByEmail($recipient['emailAddress'])->first();

                if (!empty($subscriber)) {
                    $subscriberService->updateSubscriber(['blacklisted' => true], $subscriber->id);
                } 
            }
        }
    }

    public function unsubscribe(Request $request)
    {
        //
    }
}
