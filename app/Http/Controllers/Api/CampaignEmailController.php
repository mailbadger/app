<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;
use newsletters\Http\Requests;
use newsletters\Http\Controllers\Controller;
use newsletters\Services\EmailService;

class CampaignEmailController extends Controller
{
    /**
     * @var EmailService 
     */
    private $service;

    public function __construct(EmailService $service)
    { 
        $this->service = $service;
    }

    /**
     * Display a listing of the resource.
     *
     * @return \Illuminate\Http\Response
     */
    public function index($campaignId, Request $request)
    {
        $emails = $this->service->findAllSentEmailsByCampaignId($campaignId, $request->has('paginate'), 10);     
        
        $complaints = $this->service->findSendsReportByCampaignId($campaignId);
        dd($complaints); 
        return response()->json($emails, 200);
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
}
