<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;
use newsletters\Http\Requests;
use newsletters\Http\Controllers\Controller;
use Illuminate\Support\Facades\Log;
use Aws\Sns\Message;
use Aws\Sns\MessageValidator;
use GuzzleHttp\Client;

class EmailController extends Controller
{
    /**
     * Display a listing of the resource.
     *
     * @return \Illuminate\Http\Response
     */
    public function index()
    {
        //
    }

    /**
     * Store a newly created resource in storage.
     *
     * @param  \Illuminate\Http\Request  $request
     * @return \Illuminate\Http\Response
     */
    public function store(Request $request)
    {
        //
    }

    /**
     * Display the specified resource.
     *
     * @param  int  $id
     * @return \Illuminate\Http\Response
     */
    public function show($id)
    {
        //
    }

    /**
     * Update the specified resource in storage.
     *
     * @param  \Illuminate\Http\Request  $request
     * @param  int  $id
     * @return \Illuminate\Http\Response
     */
    public function update(Request $request, $id)
    {
        //
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int  $id
     * @return \Illuminate\Http\Response
     */
    public function destroy($id)
    {
        //
    }

    public function bounces(MessageValidator $validator)
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
            Log::info('Amazon SNS bounce data: ' . print_r($message, true)); 
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
            Log::info('Amazon SNS complaints data: ' . print_r($message, true)); 
        }
    }
}
